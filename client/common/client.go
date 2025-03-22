package common

import (
	"fmt"
	"bufio"
	"os"
	"net"
	"time"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/communication"

	"github.com/op/go-logging"
)

const ReadTimeout = 1 * time.Second

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	done   chan struct{}
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		done: make(chan struct{}),
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
		return err
	}
	c.conn = conn
	return nil
}

// Shutdown Gracefully shutdown the client
func (c *Client) Shutdown() {
	log.Infof("action: graceful_shutdown | result: in_progress | client_id: %v", c.config.ID)
	close(c.done)
	c.CloseConnection()
	log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
}

// CloseConnection Closes the client connection
func (c *Client) CloseConnection() {
	if c.conn != nil {
		log.Infof("action: close_connection | result: in_progress | client_id: %s", c.config.ID)
		c.conn.Close()
		c.conn = nil
		log.Infof("action: close_connection | result: success | client_id: %s", c.config.ID)
	}
}

// ReceiveConfirmation Receives the confirmation message
func (c *Client) ReceiveConfirmation() (communication.ConfirmationMessage, error) {
	serializedPayload, err := communication.ReadMessagePayload(c.conn)
	if err != nil {
		return communication.ConfirmationMessage{}, err
	}
	
	return communication.DeserializeConfirmation(serializedPayload), nil
}

// HandleConfirmation Handles the confirmation message
func (c *Client) HandleConfirmation(confirmationMessage communication.ConfirmationMessage) {
	if confirmationMessage.Status == "success" {
		log.Infof("action: batch_apuestas_enviado | result: success")
	} else {
		log.Infof("action: batch_apuestas_enviado | result: fail")
	}
}

// ReadBetsInBatches Reads bets from CSV file in batches
func (c *Client) ReadBetsInBatches() ([][]communication.BetMessage, error) {
    filePath := fmt.Sprintf("/.data/agency-%s.csv", c.config.ID)
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("error al abrir archivo CSV: %v", err)
    }
    defer file.Close()
    
    var batches [][]communication.BetMessage
    currentBatch := []communication.BetMessage{}
    
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Split(line, ",")
        
        if len(fields) != 5 {
            return nil, fmt.Errorf("la apuesta debe tener 5 campos")
        }

        bet, err := communication.NewBetMessage(
			c.config.ID,
			fields[0],
			fields[1],
			fields[2],
			fields[3],
			fields[4],
		)
		if err != nil {
			return nil, err
		}
        
        if len(currentBatch) == c.config.BatchMaxAmount {
            batches = append(batches, currentBatch)
            currentBatch = []communication.BetMessage{}
        }
        currentBatch = append(currentBatch, bet)
    }
    
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error al leer archivo: %v", err)
    }
    
    if len(currentBatch) > 0 {
        batches = append(batches, currentBatch)
    }
    
    return batches, nil
}

// StartClientLoop Runs the client
func (c *Client) StartClientLoop() {
	batches, err := c.ReadBetsInBatches()
	if err != nil {
		log.Errorf("action: read_bets | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return
	}

	amountBatches := len(batches)

	for i, batch := range batches {
		if c.createClientSocket() != nil {
			c.CloseConnection()
			return
		}

		serializedBatch, err := communication.SerializeBatch(batch)
		if err != nil {
			log.Errorf("action: serialize_batch | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			c.CloseConnection()
			return
		}
		if err := communication.SendMessage(c.conn, serializedBatch); err != nil {
			log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			c.CloseConnection()
			return
		}

		confirmationMessage, err := c.ReceiveConfirmation()
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			c.CloseConnection()
			return
		}

		c.HandleConfirmation(confirmationMessage)
		c.conn.Close()

		if i == amountBatches - 1 {
			break
		}
		
		select {
		case <-c.done:
			return
		case <-time.After(c.config.LoopPeriod):
		}
	}
}
