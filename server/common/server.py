import os
import socket
import signal
import logging
import threading

from common.client_bet_socket import Client_bet_socket
from common.utils import store_bets,get_winners

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.is_shutting_down = False
        self.agencies_sockets = []
        self.clients_threads = []
        self.bets_utils_lock = threading.Lock()  # lock for bets load

        signal.signal(signal.SIGTERM, self.handle_sigterm)

    def handle_sigterm(self, signum, frame):
        logging.info("action: shutting_down | result: in_progress")
        self.is_shutting_down = True       

        self.gracefull_shutdown()
        
        logging.info("action: shutting_down | result: success")    
        exit(0)

    def gracefull_shutdown(self):
        for thread in self.clients_threads:
            thread.join()
        for agency_socket in self.agencies_sockets:
            try:
                # Check if socket is still valid before calling getpeername
                peer = agency_socket.socket.getpeername()
                agency_socket.close()
                logging.info(f'action: clossing_client_socket | result: success | ip: {peer}')
            except OSError:
                logging.info('action: clossing_client_socket | result: socket_already_closed')

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        client_amount = int(os.getenv("CLIENT_AMOUNT", "0"))
        barrier = threading.Barrier(client_amount + 1)

        try: 
            while not self.is_shutting_down and len(self.clients_threads) < client_amount:# tambien modif en el ej7

                new_client_socket = self.__accept_new_connection()
                if self.is_shutting_down:
                    break
                
                client_bet_socket = Client_bet_socket(new_client_socket)
                self.agencies_sockets.append(client_bet_socket)

                client_thread = threading.Thread(target=self.handle_client_connection, args=(client_bet_socket,barrier,))
                self.clients_threads.append(client_thread)
                client_thread.start()

            barrier.wait()
            logging.info('action: sorteo | result: success')
            self.handle_results()   

        finally:
            self.gracefull_shutdown()
    
    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')

        try:
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')

        except OSError as e:
            if self.is_shutting_down:
                return
            logging.error(f"action: accept_connections | result: fail | error: {e}")

        return c

    def handle_client_connection(self, client_bet_socket: Client_bet_socket, barrier: threading.Barrier):
        bets = client_bet_socket.handle_client_connection()

        with self.bets_utils_lock:
            store_bets(bets)

        barrier.wait()
                

    def handle_results(self):
        winners = get_winners()

        for agency_socket in self.agencies_sockets:

            agency_id = agency_socket.handle_winner_rqst()
            agency_winners = []
            if agency_id in winners:
                agency_winners = winners[agency_id]
            
            client_thread = threading.Thread(target=agency_socket.notif_winners, args=(agency_winners,))
            self.clients_threads.append(client_thread)
            client_thread.start()
            logging.info(f'action: send_winners | result: success | agency: {agency_id} | cant: {len(agency_winners)}')