from datetime import date
import logging
import struct

from server.common.betData import BetData


class Client_bet_socket:
    def __init__(self, socket):
        # Initialize server socket
        self.socket = socket

    def handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            bet = self.get_client_bet()
            addr = self.socket.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            self.socket.send("{}\n".format("Ok").encode('utf-8'))
        except OSError as e:
            if self.is_shutting_down:
                return
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            self.socket.close()
            logging.info(f'action: clossing_client_socket | result: success | ip: addr[0]')
        
        return bet

    def recv_exact(self, n: int):
        """Read exactly n bytes from socket, handling partial reads."""
        data = b''
        while len(data) < n:
            chunk = self.sock.recv(n - len(data))
            if not chunk:
                raise ConnectionError(f"Socket closed before reading {n} bytes")
            data += chunk
        return data


    def deserialize_bet_data_from_socket(self):
        # Agency (1 byte)
        agency = struct.unpack('>B', self.recv_exact(1))[0]

        # Document (4 bytes)
        document = struct.unpack('>I', self.recv_exact(4))[0]

        # BetNumber (8 bytes)
        bet_number = struct.unpack('>Q', self.recv_exact(8))[0]

        # Birthdate (2 + 1 + 1 bytes)
        year  = struct.unpack('>H', self.recv_exact(2))[0]
        month = struct.unpack('>B', self.recv_exact(1))[0]
        day   = struct.unpack('>B', self.recv_exact(1))[0]
        birthdate = date(year, month, day)

        # FirstName (1 byte length + N bytes)
        first_name_len = struct.unpack('>B', self.recv_exact(1))[0]
        first_name = self.recv_exact(first_name_len).decode('utf-8')

        # LastName (1 byte length + N bytes)
        last_name_len = struct.unpack('>B', self.recv_exact(1))[0]
        last_name = self.recv_exact(last_name_len).decode('utf-8')

        return BetData(agency,document,bet_number,birthdate,first_name,last_name)
