import os
import socket
import signal
import logging

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
        signal.signal(signal.SIGTERM, self.handle_sigterm)

    def handle_sigterm(self, signum, frame):
        logging.info("action: shutting_down | result: in_progress")
        self.is_shutting_down = True       

        self.gracefull_shutdown()
        
        logging.info("action: shutting_down | result: success")    
        exit(0)

    def gracefull_shutdown(self):
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

        # the server
        agencies_done = 0
        
        client_amount = int(os.getenv("CLIENT_AMOUNT", "0"))
        try: 
            while not self.is_shutting_down:
                if agencies_done >= client_amount:
                    break
                new_client_socket = self.__accept_new_connection()
                
                if self.is_shutting_down:
                    break
                
                client_bet_socket = Client_bet_socket(new_client_socket)
                bets = client_bet_socket.handle_client_connection() # NO CERRAR EL SOCKET
                store_bets(bets)
                self.agencies_sockets.append(client_bet_socket)
                agencies_done += 1

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

    def handle_results(self):
        winners = get_winners()

        for agency_socket in self.agencies_sockets:
            agency_id = agency_socket.handle_winner_rqst()
            logging.info(f'action: pasando winners a agency: {agency_id}')
            #logging.info(f'action winners: {vars(winners)}')
            agency_winners = []
            if agency_id in winners:
                agency_winners = winners[agency_id]
            
            agency_socket.notif_winners(agency_winners)
            logging.info(f'action: send_winners | result: success | agency: {agency_id} | cant: {len(agency_winners)}')