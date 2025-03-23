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

// ReceiveMessage receives a message through the connection
func ReceiveMessage(conn net.Conn, length uint32) ([]byte, error) {
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

// ReadMessagePayload reads the payload of a message from the connection
func ReadMessagePayload(conn net.Conn) ([]byte, error) {
    lengthBytes, err := ReceiveMessage(conn, LengthHeaderSize)
    if err != nil {
        return nil, err
    }

    var length uint32
    buf := bytes.NewReader(lengthBytes)
    if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
        return nil, err
    }

    payloadBytes, err := ReceiveMessage(conn, length)
    if err != nil {
        return nil, err
    }
    
    return payloadBytes, nil
}
