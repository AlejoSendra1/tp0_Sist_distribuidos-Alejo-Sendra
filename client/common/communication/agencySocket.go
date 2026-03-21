package communication

import (
	"fmt"
	"io"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/domain"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/vendor/github.com/op/go-logging"
)

type AgencySocket struct {
	conn   net.Conn
}

var log = logging.MustGetLogger("log")

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func CreateAgencySocket(serverAddr string, id string) *AgencySocket {
	conn, err := net.Dial("tcp", serverAddr)

	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
	}
	
	agencySocket := &AgencySocket{
		conn: conn,
	}
	return agencySocket
}

func (as *AgencySocket) SendBet(betData *domain.BetData, id string) (string,error) {
	serialized := SerializeBetData(betData)

	err := as.WriteExact(serialized,len(serialized))
	if err != nil {
		log.Errorf("action: send_bet | result: fail | client_id: %v | error: %v", id, err)
		return "", err
	}

	servResponse, err := as.GetServerResponse()

	as.conn.Close()
	log.Infof("action: socket_closed | result: success | client_id: %v",
		id,
	)

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			id,
			err,
		)
		return "", err
	} 

	return servResponse, nil
}

func (as *AgencySocket) WriteExact(content []byte, bytesToWrite int) error {
	bytesWritten := 0
	for bytesWritten < bytesToWrite {
		bytesWrittenNow, err := as.conn.Write(content)
		bytesWritten += bytesWrittenNow
		if err != nil {
			return err
		}
	}
	return nil
}

func readExact(conn net.Conn, n int) ([]byte, error) {
    buf := make([]byte, n)
    _, err := io.ReadFull(conn, buf)
    if err != nil {
        return nil, fmt.Errorf("failed to read %d bytes: %w", err)
    }
    return buf, nil
}

func (as *AgencySocket) GetServerResponse() (string,error) {
	responseLen , err := readExact(as.conn, 1)
	if err != nil {
		return "",err
	}
	
	serverAnswer , err := readExact(as.conn, int(uint8(responseLen[0])))
	if err != nil {
		return "",err
	}

	return string(serverAnswer), nil
}