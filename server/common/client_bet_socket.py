from datetime import date
import logging
import struct

from common.utils import Bet

BATCH_HEADER_SIZE = 3
WINNERS_RQST_CODE = 2

class Client_bet_socket:
    def __init__(self, sockett):
        self.socket = sockett

    def close(self):
        try:
            addr = self.socket.getpeername()
            self.socket.close()
            logging.info(f'action: clossing_client_socket | result: success | ip: {addr[0]}')
        finally:
            return

    def send_response(self,content: str):
        """Writes exactly the 'content' (string) given in the socket"""
        content_chunk = content.encode('utf-8')
        body_size = len(content_chunk).to_bytes(1,'big')
        message = body_size + content_chunk
        self.send_exact(message)
    
    def send_exact(self, bytes):
        """Writes exactly the bytes given in the socket"""
        sent = 0
        while sent < len(bytes):
            bytes_sent =  self.socket.send(bytes[sent:])
            sent += bytes_sent
    
    def handle_client_connection(self) -> list:
        """Read message from a specific client socket and closes the socket
        If a problem arises in the communication with the client, the
        client socket will also be closed """

        bets = []

        try:
            bets = self.get_client_bets()
        except OSError as e:
            return []
        except Exception as err:
            self.send_response(f'Error: {err}')
            return []
        
        return bets
        


    def recv_exact(self, n: int) -> int:
        """Reads exactly n bytes from the socket"""
        readed = b''
        while len(readed) < n:
            chunk = self.socket.recv(n - len(readed))
            if not chunk:
                raise ConnectionError(f"Socket closed before reading {n} bytes")
            readed += chunk
        return readed

    def get_client_bets(self) -> list:
        """ el primer byte/bytes corresponderan a la cantidad de bets a leer"""
        bets = []
        bytes_to_read = BATCH_HEADER_SIZE

        
        at_least_one_with_err = False
        bytes_to_read = struct.unpack('>H', self.recv_exact(2))[0]
        agency = struct.unpack('>B', self.recv_exact(1))[0]

        bets_obtained = 0
        read_bytes = 0
        while read_bytes < bytes_to_read:

            new_bet, bytes_read = self.get_client_bet(agency)
            read_bytes += bytes_read
            if new_bet is not None:
                bets_obtained += 1
                bets.append(new_bet)
            else:
                at_least_one_with_err = True
        
        if at_least_one_with_err:
            logging.info(f'action: apuesta_recibida | result: fail | cantidad: {bets_obtained}')
            self.send_response("Error: At least one bet from the batch has a wrong format")
        else:
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {bets_obtained}')
            self.send_response("Ok")

        return bets
        
        

    def get_client_bet(self, agency: int) -> (Bet,int):
        """Reads from the socket every one of the fields of the Bet Object
        respecting the protocol and return the a Bet object and the bytes read"""
        amount_bytes_read = 0

        bytes = self.recv_exact(4)
        amount_bytes_read += 4
        document = struct.unpack('>I', bytes)[0]

        bytes = self.recv_exact(8)
        amount_bytes_read += 8
        bet_number = struct.unpack('>Q', bytes)[0]

        bytes = self.recv_exact(2)
        amount_bytes_read += 2
        year  = struct.unpack('>H', bytes)[0]

        bytes = self.recv_exact(1)
        amount_bytes_read += 1
        month = struct.unpack('>B', bytes)[0]
        
        bytes = self.recv_exact(1)
        amount_bytes_read += 1
        day   = struct.unpack('>B', bytes)[0]

        bytes = self.recv_exact(1)
        amount_bytes_read += 1
        first_name_len = struct.unpack('>B', bytes)[0]

        bytes = self.recv_exact(first_name_len)
        amount_bytes_read += first_name_len
        first_name_chunks = bytes

        bytes = self.recv_exact(1)
        amount_bytes_read += 1
        last_name_len = struct.unpack('>B', bytes)[0]

        bytes = self.recv_exact(last_name_len)
        amount_bytes_read += last_name_len
        last_name_chunks = bytes

        try: 
            first_name = first_name_chunks.decode('utf-8')
            last_name = last_name_chunks.decode('utf-8')
            birthdate = date(year, month, day).isoformat()
            return Bet(agency,first_name,last_name,document,birthdate,bet_number), amount_bytes_read
        except Exception as err:
            logging.info(f'Error reading bet: {err}') # cambiar a log tipo error
            return None, amount_bytes_read
        
    def handle_winner_rqst(self): # necesito recibir el header y el numero de agency y devolver el agency
        """Reads the request msg and return the agency number to identify the agency"""
        msg_type = struct.unpack('>B', self.recv_exact(1))[0]
        agency_id = struct.unpack('>B', self.recv_exact(1))[0]
        
        if msg_type != WINNERS_RQST_CODE:
            logging.error(f'action: request_for_winner | result: fail')
            return
        logging.info(f'action: notify_done | result: success | agency: {agency_id}')
        return agency_id

    def notif_winners(self,agency_winners: list): # serializa los datos y se lo manda a la agencia 
        """TO DO description"""
        bytes_to_read = struct.pack('>B', len(agency_winners))

        to_send = bytes_to_read
        for document in agency_winners:
            to_send += struct.pack('>I', int(document))
        
        self.send_exact(to_send)