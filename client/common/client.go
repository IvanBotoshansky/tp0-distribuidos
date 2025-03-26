package common

import (
	"fmt"
	"bufio"
	"os"
	"net"
	"time"
	"strings"
	"sync"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/communication"

	"github.com/op/go-logging"
)

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
	connMutex sync.Mutex
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
	c.closeConnection()
	close(c.done)
	log.Infof("action: graceful_shutdown | result: success | client_id: %v", c.config.ID)
}

// closeConnection Closes the client connection
func (c *Client) closeConnection() {
	c.connMutex.Lock()
	defer c.connMutex.Unlock()
	if c.conn != nil {
		log.Infof("action: close_connection | result: in_progress | client_id: %s", c.config.ID)
		c.conn.Close()
		c.conn = nil
		log.Infof("action: close_connection | result: success | client_id: %s", c.config.ID)
	}
}

// receiveConfirmation Receives the confirmation message
func (c *Client) receiveConfirmation() (communication.ConfirmationMessage, error) {
	messageType, serializedPayload, err := communication.ReceiveMessage(c.conn)
	if err != nil {
		return communication.ConfirmationMessage{}, err
	}
	if messageType != communication.MessageTypeConfirmation {
		return communication.ConfirmationMessage{}, fmt.Errorf("mensaje inesperado")
	}
	return communication.DeserializeConfirmation(serializedPayload), nil
}

// handleConfirmation Handles the confirmation message
func (c *Client) handleConfirmation(confirmationMessage communication.ConfirmationMessage) {
	if confirmationMessage.Status == "success" {
		log.Infof("action: batch_apuestas_enviado | result: success")
	} else {
		log.Infof("action: batch_apuestas_enviado | result: fail")
	}
}

// readBetsInBatches Reads bets from CSV file in batches
func (c *Client) readBetsInBatches() ([][]communication.BetMessage, error) {
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

// sendBatch serializes and sends a batch of bets
func (c *Client) sendBatch(batch []communication.BetMessage) error {
    serializedBatch, err := communication.SerializeBatch(batch)
	if err != nil {
		return err
	}
	if err := communication.SendMessage(c.conn, serializedBatch); err != nil {
		return err
	}
	return nil
}

// sendEndNotification serializes and sends an end notification message
func (c *Client) sendEndNotification() error {
	endNotification := communication.NewEndNotificationMessage(c.config.ID)
	serializedEndNotification, err := communication.SerializeEndNotification(endNotification)
	if err != nil {
		return err
	}
	if err := communication.SendMessage(c.conn, serializedEndNotification); err != nil {
		return err
	}
	return nil
}

// sendWinnersRequest serializes and sends a winners request message
func (c *Client) sendWinnersRequest() error {
	winnersRequest := communication.NewWinnersRequestMessage(c.config.ID)
	serializedWinnersRequest, err := communication.SerializeWinnersRequest(winnersRequest)
	if err != nil {
		return err
	}
	if err := communication.SendMessage(c.conn, serializedWinnersRequest); err != nil {
		return err
	}
	return nil
}

// sendAllBatches Sends all batches of bets
func (c *Client) sendAllBatches(batches [][]communication.BetMessage) bool {
	amountBatches := len(batches)

	for i, batch := range batches {
		if err := c.sendBatch(batch); err != nil {
			log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
			return false
		}

		confirmationMessage, err := c.receiveConfirmation()
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID, err)
			return false
		}

		c.handleConfirmation(confirmationMessage)

		if i == amountBatches - 1 {
			break
		}
	}
	return true
}

// notifyEndOfBets Notifies the end of bets sending
func (c *Client) notifyEndOfBets() bool {
	if err := c.sendEndNotification(); err != nil {
		log.Errorf("action: send_end_notification | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false
	}
	return true
}

// askForWinners Asks for winners
func (c *Client) askForWinners() bool {
	if err := c.sendWinnersRequest(); err != nil {
		log.Errorf("action: send_winners_request | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false
	}

	messageType, serializedPayload, err := communication.ReceiveMessage(c.conn)
	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return false
	}
	if messageType == communication.MessageTypeWinnersList {
		winnersList := communication.DeserializeWinnersList(serializedPayload)
		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winnersList.Winners))
	} else {
		log.Errorf("action: consulta_ganadores | result: fail")
		return false
	}

	return true
}

// StartClientLoop Runs the client
func (c *Client) StartClientLoop() {
	batches, err := c.readBetsInBatches()
	if err != nil {
		log.Errorf("action: read_bets | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		return
	}

	if err := c.createClientSocket(); err != nil {
		log.Errorf("action: initial_connection | result: fail | client_id: %v | error: %v",
            c.config.ID, err)
		return
	}

	defer c.closeConnection()

	if !c.sendAllBatches(batches) {
		return
	}

	if !c.notifyEndOfBets() {
		return
	}

	if !c.askForWinners() {
		return
	}
}
