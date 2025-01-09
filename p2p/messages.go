package p2p

import "encoding/json"

//structured message
type Message struct {
	Type         string      `json:"type"`    //e.g., "REQUEST_CHAIN", "NEW_BLOCK"
	Dataset      interface{} `json:"dataset"` //CID of dataset to be used with the algorithm
	Algo         interface{} `json:"algo"`    //CID of algo to be used
	Requirements interface{} `json:"req"`     //CID requirements file to be installed
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
