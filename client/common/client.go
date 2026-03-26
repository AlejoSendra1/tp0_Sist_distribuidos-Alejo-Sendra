package common

import (
	"os"
	"os/signal"
	"syscall"

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
	BatchAmount    int
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
func (c *Client) StartClient() error {

	// Create the connection to the server
	agencySocket, err:= communication.CreateAgencySocket(c.config.ServerAddress,c.config.ID)
	if err != nil {
		log.Criticalf("%s", err)
		return err
	}

	log.Infof("action: apuestas_enviadas | result: in_progress | client_id: %v", c.config.ID)
	
	reader, err := domain.NewReader(c.config.ID,c.config.BatchAmount)
	defer reader.Close()
	if err != nil {
		return err
    }
	log.Infof("action: apuestas_enviadas | result: in_progress")

	for {
		log.Infof("action: get_data_csv | result: in_progress")
		bets, dataErr := reader.GetBatch()
		if dataErr != nil {
			log.Infof("action: get_data_csv | result: fail")
			log.Criticalf("%v", err)
			return dataErr
		}
		log.Infof("action: get_data_csv | result: success")
		
		err = agencySocket.SendBets(bets,c.config.ID) 
		if err != nil {
			log.Criticalf("%s", err)
			return err
		}

		if len(bets) == 0 {
			break
		}
	}
	
	agencySocket.GetWinners()
	
	log.Infof("action: shutting_down_default | result: in_progress | client_id: %v", c.config.ID)
    select {
    case <-c.signalChannel:
        log.Infof("action: SIGTERM_caught | result: success | client_id: %v", c.config.ID)
        agencySocket.Close() 
        log.Infof("action: shutting_down | result: success | client_id: %v", c.config.ID)
        return nil
	default:
		agencySocket.Close() 
		log.Infof("action: shutting_down_default | result: success | client_id: %v", c.config.ID)
		os.Exit(0)
    }
	return nil
}