package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/domain"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	betData domain.BetData
	conn   net.Conn
	signalChannel chan os.Signal
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, clientBetData domain.BetData) *Client {

	channel:= make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGTERM, os.Interrupt)

	client := &Client{
		config: config,
		betData: clientBetData,
		signalChannel: channel,
	}

	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {

	//betSerialization := serialization.serializeBet(c.bet)

	// Create the connection the server
	c.createClientSocket()

	// TODO: Modify the send to avoid short-write
	log.Infof("action: apuesta_enviada | result: in_progress | dni: %v | numero: %v | agency: %v",
		c.betData.Document,
		c.betData.BetNumber,
		c.config.ID,
	)
	fmt.Fprintf(
		c.conn,
		"", // aca enviar bytess/serializacion 
	)
	serverAnswer, err := bufio.NewReader(c.conn).ReadString('\n')
	
	c.conn.Close()
	log.Infof("action: socket_closed | result: success | client_id: %v",
		c.config.ID,
	)

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	} else {
		log.Infof("action: msg_recivido del server %s",
			serverAnswer,
		)
	}

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v | agency: %v",
		c.betData.Document,
		c.betData.BetNumber,
		c.config.ID,
	)

	select {
		case <-c.signalChannel:
			log.Infof("action: SIGTERM_caught | result: success | client_id: %v", c.config.ID)
			log.Infof("action: shutting_down | result: success | client_id: %v", c.config.ID)
		return
		case <-time.After(c.config.LoopPeriod):
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
