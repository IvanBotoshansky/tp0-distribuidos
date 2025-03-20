package common

import (
	"net"
	"time"

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
	FirstName     string
	LastName      string
	Document      string
	Birthdate     string
	Number        string
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
	log.Infof("action: close_connection | result: in_progress | client_id: %s", c.config.ID)
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	log.Infof("action: close_connection | result: success | client_id: %s", c.config.ID)
}

// ReadPayloadWithTimeout Reads message payload with a timeout
func (c *Client) ReadPayloadWithTimeout() ([]byte, error) {
    if c.conn != nil {
        _ = c.conn.SetReadDeadline(time.Now().Add(ReadTimeout))
    }
    return communication.ReadMessagePayload(c.conn)
}

// IsTimeoutError Checks if the error is a timeout error
func IsTimeoutError(err error) bool {
    netErr, ok := err.(net.Error)
    return ok && netErr.Timeout()
}

// ReceiveConfirmation Receives the confirmation message
func (c *Client) ReceiveConfirmation(confirmationc chan communication.ConfirmationMessage, errorc chan error) {
	for {
		serializedPayload, err := c.ReadPayloadWithTimeout()
		select {
		case <-c.done:
			return
		default:
		}

		if err != nil {
			if IsTimeoutError(err) {
				continue
			}
			errorc <- err
			return
		}
		
		confirmationMessage := communication.DeserializeConfirmation(serializedPayload)
		confirmationc <- confirmationMessage
		return
	}
}

// HandleConfirmation Handles the confirmation message
func (c *Client) HandleConfirmation(confirmationMessage communication.ConfirmationMessage) {
	if confirmationMessage.Status == "success" {
		log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s", c.config.Document, c.config.Number)
	} else {
		log.Infof("action: apuesta_enviada | result: fail | dni: %s | numero: %s", c.config.Document, c.config.Number)
	}
}

// StartClientLoop Runs the client
func (c *Client) StartClientLoop() {
	if c.createClientSocket() != nil {
		c.CloseConnection()
		return
	}

	betMessage := communication.NewBetMessage(c.config.ID, c.config.FirstName, c.config.LastName, c.config.Document, c.config.Birthdate, c.config.Number)
	serializedBet, err := communication.SerializeBet(betMessage)
	if err != nil {
		log.Errorf("action: serialize_bet | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		c.CloseConnection()
		return
	}
	if err := communication.SendMessage(c.conn, serializedBet); err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
			c.config.ID, err)
		c.CloseConnection()
		return
	}
	
	confirmationc := make(chan communication.ConfirmationMessage)
	errorc := make(chan error)

	go c.ReceiveConfirmation(confirmationc, errorc)

	select {
    case confirmationMessage := <-confirmationc:
		c.HandleConfirmation(confirmationMessage)
		c.CloseConnection()
	case err := <-errorc:
        log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
            c.config.ID, err)
		c.CloseConnection()
	case <-c.done:
	}
}
