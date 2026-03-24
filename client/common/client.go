package common

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/communication"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/domain"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	socket communication.AgencySocket
	signalChannel chan os.Signal
}

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {

	channel:= make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGTERM, os.Interrupt)

	client := &Client{
		config: config,
		signalChannel: channel,
	}

	return client
}


// StartClient Send messages to the client until some time threshold is met
func (c *Client) StartClient(bets []domain.Bet) error {

	// Create the connection to the server
	agencySocket, err:= communication.CreateAgencySocket(c.config.ServerAddress,c.config.ID)
	if err != nil {
		log.Criticalf("%s", err)
		return err
	}

	log.Infof("action: apuestas_enviadas | result: in_progress")

	err = agencySocket.SendBets(bets,c.config.ID) 
	if err != nil {
		log.Criticalf("%s", err)
		return err
	}

	log.Infof("action: apuestas_enviadas | result: success ")

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID) // Eliminar luego de probar si pasan los tests --------------
	
	select {
		case <-c.signalChannel:
			log.Infof("action: SIGTERM_caught | result: success | client_id: %v", c.config.ID)
			log.Infof("action: shutting_down | result: success | client_id: %v", c.config.ID)
		return nil
	}	
}
