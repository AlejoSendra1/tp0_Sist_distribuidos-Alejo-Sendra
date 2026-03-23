package communication

import (
	"io"
	"net"

	"github.com/op/go-logging"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/domain"
)

type AgencySocket struct {
	conn   net.Conn
}

var log = logging.MustGetLogger("log")

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
	}

	log.Infof(
		"action: connect | result: success | client_id: %v",
		id,
	)

	return agencySocket, nil
}

func (as *AgencySocket) SendBet(betData domain.Bet, id string) (string,error) {
	// Sends the given Bet to the server and returns the server response 
	// in case communication is successfull
	serialized := SerializeBet(&betData)

	err := as.writeExact(serialized)
	if err != nil {
		log.Criticalf("action: send_bet | result: fail | client_id: %v | error: %v", id, err)
		return "", err
	}

	servResponse, err := as.GetServerResponse()

	as.conn.Close()
	log.Infof("action: socket_closed | result: success | client_id: %v",
		id,
	)

	if err != nil {
		log.Criticalf("action: receive_message | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
		return "", err
	} 

	return servResponse, nil
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
		log.Infof("action: sent | result: success | content: %v | amount_of_bytes: %v",
			content[bytesWritten:],
			len(content),
		)
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