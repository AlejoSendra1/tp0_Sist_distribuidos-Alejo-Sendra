import socket
import signal
import logging

from server.common.communication_utils import Client_bet_socket


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.is_shutting_down = False
        signal.signal(signal.SIGTERM, self.handle_sigterm)

    def handle_sigterm(self, signum, frame):
        logging.info("action: shutting_down | result: in_progress")
        self.is_shutting_down = True       

        try:
            logging.info("action: clossing_listening_socket | result: in_progress")
            self._server_socket.close()
            logging.info("action: clossing_listening_socket | result: success")
        except OSError as e:
            logging.info("action: clossing_listening_socket | result: fail")
            exit(1)

        try:
            if self.client_sock is not None:
                logging.info("action: clossing_client_socket | result: in_progress")
                self.client_sock.close()
                logging.info("action: clossing_client_socket | result: success")
        except OSError as e:
            logging.info("action: clossing_client_socket | result: fail")
            exit(1)
        
        logging.info("action: shutting_down | result: success")    
        exit(0)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # the server
        while not self.is_shutting_down:
            new_client_socket = self.__accept_new_connection()
            
            if self.is_shutting_down:
                break
            
            client_bet_socket = Client_bet_socket(new_client_socket)#crear class
            client_bet_socket.handle_client_connection()  # obtener la apuesta
            # guardarla 
            

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
            logging.error("action: accept_connections | result: fail | error: {e}")

        return c
