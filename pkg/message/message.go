/*
 * Revision History:
 *     Initial: 2018/05/25        Tong Yuehong
 */

package message

import (
	json "github.com/json-iterator/go"

	"github.com/TechCatsLab/rumour"
	"github.com/TechCatsLab/rumour/identify"
)

type Message struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// Type return message's type.
func (msg *Message) MessageType() interface{} {
	return msg.Type
}

// Source return the identify which send the message.
func (msg *Message) Source() rumour.Identify {
	return identify.Identify(msg.From)
}

// Target return the identify which receive the message.
func (msg *Message) Target() rumour.Identify {
	return identify.Identify(msg.To)
}

// Marshal -
func (msg *Message) Marshal() ([]byte, error) {
	message, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Unmarshal -
func (msg *Message) Unmarshal(info []byte) error {
	return json.Unmarshal(info, &msg)
}
