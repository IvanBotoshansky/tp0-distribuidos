package communication

import (
    "bytes"
    "encoding/binary"
    "strings"
)

const BetFieldsSeparator = ","

// SerializeData serializes data
func SerializeData(data string) ([]byte, error) {
    dataBytes := []byte(data)
    length := uint32(len(dataBytes))
    buf := new(bytes.Buffer)
    
    if err := binary.Write(buf, binary.BigEndian, length); err != nil {
        return nil, err
    }
    
    if _, err := buf.Write(dataBytes); err != nil {
        return nil, err
    }
    
    return buf.Bytes(), nil
}

// SerializeBet serializes a bet message
func SerializeBet(bet BetMessage) ([]byte, error) {
    data := strings.Join([]string{
        bet.Agency, 
        bet.FirstName, 
        bet.LastName, 
        bet.Document, 
        bet.Birthdate, 
        bet.Number,
    }, BetFieldsSeparator)
    return SerializeData(data)
}

// DeserializeConfirmation deserializes data bytes of a confirmation message
func DeserializeConfirmation(dataBytes []byte) ConfirmationMessage {
	return NewConfirmationMessage(string(dataBytes))
}
