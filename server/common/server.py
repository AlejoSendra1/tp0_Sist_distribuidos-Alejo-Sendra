import socket
import signal
import logging


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.client_sock = None
        self.is_shutting_down = False
        signal.signal(signal.SIGTERM, self.handle_sigterm)

    def handle_sigterm(self):
        logging.info("action: shutting_down | result: in_process | server")
        self.is_shutting_down = True

        logging.info("action: clossing_lisening_socket | server")
        self._server_socket.close()
        logging.info("action: clossing_lisening_socket | result: success | server")

        logging.info("action: clossing_client_socket | server")
        if self.client_sock:
            self.client_sock.close()
        logging.info("action: clossing_client_socket | result: success | server")

        logging.info("action: shutting_down | result: success | server")
        


    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # the server
        while True:
            self.__accept_new_connection()
            self.__handle_client_connection()

    def __handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = self.client_sock.recv(1024).rstrip().decode('utf-8')
            addr = self.client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            self.client_sock.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            if self.is_shutting_down:
                return
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            self.client_sock.close()

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
            self.client_sock = c
        except OSError as e:
            if self.is_shutting_down:
                return
            logging.error("action: accept_connections | result: fail | error: {e}")
        finally:
            self._server_socket.close()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return 
