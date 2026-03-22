from datetime import date
import logging
import struct

from common.utils import Bet


class Client_bet_socket:
    def __init__(self, socket):
        # Initialize server socket
        self.socket = socket

    def handle_client_connection(self) -> Bet:
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        addr = self.socket.getpeername()
        bet = None

        try:
            bet = self.get_client_bet()
            self.send("Ok")

        except OSError as e:
            if self.is_shutting_down:
                return
            logging.error("action: receive_message | result: fail | error: {e}")

        except Exception as err:
            self.send(f'Error: {err}')
            return
        finally:
            self.socket.close()
            logging.info(f'action: clossing_client_socket | result: success | ip: {addr[0]}')
        return bet
        


    def recv_exact(self, n: int) -> int:
        """Reads exactly n bytes from the socket"""
        readed = b''
        while len(readed) < n:
            chunk = self.socket.recv(n - len(readed))
            if not chunk:
                raise ConnectionError(f"Socket closed before reading {n} bytes")
            readed += chunk
        return readed


    def get_client_bet(self) -> Bet:
        """Reads from the socket every one of the fields of the Bet Object
        respecting the protocol"""
        
        agency = struct.unpack('>B', self.recv_exact(1))[0]
        logging.info(f'agencia: {agency}')

        document = struct.unpack('>I', self.recv_exact(4))[0]
        logging.info(f'document: {document}')

        bet_number = struct.unpack('>Q', self.recv_exact(8))[0]
        logging.info(f'bet_number: {bet_number}')

        year  = struct.unpack('>H', self.recv_exact(2))[0]
        month = struct.unpack('>B', self.recv_exact(1))[0]
        day   = struct.unpack('>B', self.recv_exact(1))[0]
        logging.info(f'year: {year}, month: {month}, day: {day}')

        first_name_len = struct.unpack('>B', self.recv_exact(1))[0]
        first_name_chunks = self.recv_exact(first_name_len)

        last_name_len = struct.unpack('>B', self.recv_exact(1))[0]
        last_name_chunks = self.recv_exact(last_name_len)

        first_name = first_name_chunks.decode('utf-8')
        logging.info(f'first_name: {first_name}')
        last_name = last_name_chunks.decode('utf-8')
        logging.info(f'last_name: {last_name}')
        birthdate = date(year, month, day).isoformat()
        return Bet(agency,first_name,last_name,document,birthdate,bet_number)
    
    def send(self,content: str):
        content_chunk = content.encode('utf-8')
        body_size = len(content_chunk).to_bytes(1,'big')

        message = body_size + content_chunk

        sent = 0
        while sent < len(message):
            bytes_sent =  self.socket.send(message[sent:])
            sent += bytes_sent

        
        


