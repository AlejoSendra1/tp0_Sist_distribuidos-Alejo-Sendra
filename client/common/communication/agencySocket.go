package communication

import (
	"io"
	"net"
	"strconv"
	"encoding/binary"
	"github.com/op/go-logging"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/domain"
)

type AgencySocket struct {
	conn   net.Conn
	id string
}

var log = logging.MustGetLogger("log")

const (
	WINNERS_RQST_CODE = byte(0x02) // TODO convertir en env var y tmb la de python
	WINNERS_RQST_SIZE = 2
)

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func CreateAgencySocket(serverAddr string, id string) (*AgencySocket, error) {
	conn, err := net.Dial("tcp", serverAddr)

	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
		return &AgencySocket{}, err
	}
	
	agencySocket := &AgencySocket{
		conn: conn,
		id: id,
	}

	log.Infof(
		"action: connect | result: success | client_id: %v",
		id,
	)

	return agencySocket, nil
}



func (as *AgencySocket) SendBets(bets []domain.Bet, id string) error {
	// Sends the given Bet to the server and returns the server response 
	// in case communication is successfull
	
	//defer as.Close()
	serialized := createBatch(bets, id)

	err := as.writeExact(serialized)
	if err != nil {
		log.Criticalf("action: send_bets | result: fail | client_id: %v | error: %v", id, err)
		return err
	}
	log.Infof("action: send_bets | result: success | se_enviaron: \"%v\" en el batch", len(bets))

	servResponse, err := as.GetServerResponse()
	if err != nil {
		log.Criticalf("action: receive_message | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
	} else {
		log.Infof("action: server_response | result: success | answer: %v",
			servResponse,
		)
	}

	return nil
}

func (as *AgencySocket) writeExact(content []byte) error {
	// Writes exactly the content given in the socket
	bytesWritten := 0
	for bytesWritten < len(content) {
		bytesWrittenNow, err := as.conn.Write(content[bytesWritten:])
		bytesWritten += bytesWrittenNow
		if err != nil {
			return err
		}
		/*
		log.Infof("action: sent | result: success | amount_of_bytes: %v",
			len(content),
		)
			*/
	}
	return nil
}

func recvExact(conn net.Conn, n int) ([]byte, error) {
	//  reads exactly n bytes from the socket
	//  and returns the array of bytes read and an error if fewer bytes were read
    buf := make([]byte, n)
    _, err := io.ReadFull(conn, buf)
    if err != nil {
		log.Criticalf("failed to read %v bytes: %v", n, err)
        return nil, err
    }
    return buf, nil
}

func (as *AgencySocket) GetServerResponse() (string,error) {
	// Wait for the server response and return its answer in case of success
	responseLen , err := recvExact(as.conn, 1)
	if err != nil {
		return "",err
	}
	
	serverAnswer , err := recvExact(as.conn, int(uint8(responseLen[0])))
	if err != nil {
		return "",err
	}

	return string(serverAnswer), nil
}

func (as *AgencySocket) Close() {
	as.conn.Close()
	log.Infof("action: socket_closed | result: success | client_id: %v",
			as.id,		
		)

}

// Sends to the server a request for the winners and waits for the answer
func (as *AgencySocket) GetWinners() error {
	// request
	rqstForWinners := make([]byte, HEADER_SIZE)
	rqstForWinners[0] = WINNERS_RQST_CODE
    agencyNumAsInt, _ := strconv.Atoi(as.id)
    rqstForWinners[1] = byte(agencyNumAsInt)
 
	if err := as.writeExact(rqstForWinners); err != nil {
		log.Criticalf("action: consulta_ganadores | result: fail | error: %v", err)
		return err
	}
	log.Infof("action: consulta_ganadores | result: in_progress | client_id: %v", as.id)
 
	winnersAmount, err := recvExact(as.conn, 1)
	if err != nil {
		log.Criticalf("action: consulta_ganadores | result: fail | error: %v", err)
		return err
	}
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", winnersAmount[0])

	winners := make([]uint32, uint8(winnersAmount[0]))
    for i := uint8(0); i < uint8(winnersAmount[0]); i++ {
        dniBytes, err := recvExact(as.conn, 4) 
        if err != nil {
            log.Criticalf("action: receive_winners | result: fail | error: %v", err)
            return err
        }
        winners[i] = binary.BigEndian.Uint32(dniBytes)
    }
    
    log.Infof("action: lista_ganadores | result: success | ganadores: %v", winners)
	return nil

}