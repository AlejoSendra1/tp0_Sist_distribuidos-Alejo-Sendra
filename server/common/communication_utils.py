from datetime import date
import logging
import struct

from server.common.utils import Bet


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
        addr = self.socket.getpeername()

        try:
            bet = self.get_client_bet()
            self.socket.send("{}\n".format("Ok").encode('utf-8'))

        except OSError as e:
            if self.is_shutting_down:
                return
            logging.error("action: receive_message | result: fail | error: {e}")

        except Exception as err:
            self.socket.send("{}\n".format("Error:{err}").encode('utf-8'))
        finally:
            self.socket.close()
            logging.info(f'action: clossing_client_socket | result: success | ip: {addr[0]}')
        
        return bet


    def recv_exact(self, n: int):
        """Reads exactly n bytes from the socket"""
        data = b''
        while len(data) < n:
            chunk = self.sock.recv(n - len(data))
            if not chunk:
                raise ConnectionError(f"Socket closed before reading {n} bytes")
            data += chunk
        return data


    def get_client_bet(self):
        agency = struct.unpack('>B', self.recv_exact(1))[0]

        document = struct.unpack('>I', self.recv_exact(4))[0]

        bet_number = struct.unpack('>Q', self.recv_exact(8))[0]

        year  = struct.unpack('>H', self.recv_exact(2))[0]
        month = struct.unpack('>B', self.recv_exact(1))[0]
        day   = struct.unpack('>B', self.recv_exact(1))[0]

        first_name_len = struct.unpack('>B', self.recv_exact(1))[0]
        first_name_chunks = self.recv_exact(first_name_len)

        last_name_len = struct.unpack('>B', self.recv_exact(1))[0]
        last_name_chunks = self.recv_exact(last_name_len)

        first_name = first_name_chunks.decode('utf-8')
        last_name = last_name_chunks.decode('utf-8')
        birthdate = date(year, month, day)
        return Bet(agency,first_name,last_name,document,birthdate,bet_number)