package communication

import (
    "bytes"
    "encoding/binary"
    "strings"
)

const (
    BetFieldsSeparator = ","
    BetsSeparator = ";"
)

const MessageTypeBetBatch        byte = 1

// SerializeWithType serializes data with a message type
func SerializeWithType(messageType uint8, data string) ([]byte, error) {
    dataBytes := []byte(data)
    length := uint32(1 + len(dataBytes))
    buf := new(bytes.Buffer)

    if err := binary.Write(buf, binary.BigEndian, length); err != nil {
        return nil, err
    }
    if err := buf.WriteByte(messageType); err != nil {
        return nil, err
    }
    if _, err := buf.Write(dataBytes); err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

// SerializeBatch serializes a batch of BetMessage
func SerializeBatch(batch []BetMessage) ([]byte, error) {
    betStrings := make([]string, len(batch))
    for i, bet := range batch {
        betStrings[i] = strings.Join([]string{
            bet.Agency, 
            bet.FirstName, 
            bet.LastName, 
            bet.Document, 
            bet.Birthdate, 
            bet.Number,
        }, BetFieldsSeparator)
    }
    data := strings.Join(betStrings, BetsSeparator)
    return SerializeWithType(MessageTypeBetBatch, data)
}

// DeserializeConfirmation deserializes data bytes of a confirmation message
func DeserializeConfirmation(dataBytes []byte) ConfirmationMessage {
	return NewConfirmationMessage(string(dataBytes))
}
