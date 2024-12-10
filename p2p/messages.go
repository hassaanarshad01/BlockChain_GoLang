package p2p

import "encoding/json"

//structured message
type Message struct {
	Type string      `json:"type"` //e.g., "REQUEST_CHAIN", "NEW_BLOCK"
	Data interface{} `json:"data"` //Data can be a string, map, or a specific struct
}

//convert to JSON
func SerializeMessage(message Message) (string, error) {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	return string(messageBytes), nil
}

//convert from JSON
func DeserializeMessage(jsonString string) (Message, error) {
	var message Message
	err := json.Unmarshal([]byte(jsonString), &message)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}
