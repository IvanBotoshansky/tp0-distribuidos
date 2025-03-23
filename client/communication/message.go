package communication

import (
    "net"
    "bytes"
    "encoding/binary"
    "fmt"
)

const LengthHeaderSize = 4
const MaxNameLength = 50

// BetMessage represents the message sent by the client to register a bet
type BetMessage struct {
    Agency    string
    FirstName string
    LastName  string
    Document  string
    Birthdate string
    Number    string
}

// ConfirmationMessage represents the confirmation message received by the client
type ConfirmationMessage struct {
    Status string
}

// EndNotificationMessage represents the message sent by the client to notify the end of the bets sending
type EndNotificationMessage struct{
    Agency string
}

// WinnersRequestMessage represents the message sent by the client to request the winners
type WinnersRequestMessage struct {
    Agency string
}

// WinnersListMessage represents the message received by the client with the winners list
type WinnersListMessage struct {
    Winners []string
}

// NewBetMessage creates a new BetMessage instance
func NewBetMessage(agency, firstName, lastName, document, birthdate, number string) (BetMessage, error) {
    if len(firstName) > MaxNameLength || len(lastName) > MaxNameLength {
        return BetMessage{}, fmt.Errorf("nombre y apellido no pueden tener más de %d caracteres", MaxNameLength)
    }
    return BetMessage{
        Agency:    agency,
        FirstName: firstName,
        LastName:  lastName,
        Document:  document,
        Birthdate: birthdate,
        Number:    number,
    }, nil
}

// NewConfirmationMessage creates a new ConfirmationMessage instance
func NewConfirmationMessage(status string) ConfirmationMessage {
    return ConfirmationMessage{
        Status: status,
    }
}

// NewEndNotificationMessage creates a new EndNotificationMessage instance
func NewEndNotificationMessage(agency string) EndNotificationMessage {
    return EndNotificationMessage{
        Agency: agency,
    }
}

// NewWinnersRequestMessage creates a new WinnersRequestMessage instance
func NewWinnersRequestMessage(agency string) WinnersRequestMessage {
    return WinnersRequestMessage{
        Agency: agency,
    }
}

// NewWinnersListMessage creates a new WinnersListMessage instance
func NewWinnersListMessage(winners []string) WinnersListMessage {
    return WinnersListMessage{
        Winners: winners,
    }
}

// SendMessage sends a message through the connection
func SendMessage(conn net.Conn, data []byte) error {
    totalBytesSent := 0
    for totalBytesSent < len(data) {
        bytesSent, err := conn.Write(data[totalBytesSent:])
        if err != nil {
            return err
        }
        totalBytesSent += bytesSent
    }
    return nil
}

// ReadMessage reads a message from the connection
func ReadMessage(conn net.Conn, length uint32) ([]byte, error) {
    totalBytesReceived := 0
    data := make([]byte, length)
    for totalBytesReceived < int(length) {
        bytesReceived, err := conn.Read(data[totalBytesReceived:])
        if err != nil {
            return nil, err
        }
        totalBytesReceived += bytesReceived
    }
    return data, nil
}

// ReceiveMessage receives the payload of a message and its type from the connection
func ReceiveMessage(conn net.Conn) (byte, []byte, error) {
    lengthBytes, err := ReadMessage(conn, LengthHeaderSize)
    if err != nil {
        return 0, nil, err
    }

    var length uint32
    buf := bytes.NewReader(lengthBytes)
    if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
        return 0, nil, err
    }

    fullMessage, err := ReadMessage(conn, length)
    if err != nil {
        return 0, nil, err
    }

    messageType := fullMessage[0]
    payloadBytes := fullMessage[1:]
    
    return messageType, payloadBytes, nil
}
