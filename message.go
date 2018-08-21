/*
 * Revision History:
 *     Initial: 2018/05/25        Tong Yuehong
 */

package rumour

import (
	json "github.com/json-iterator/go"
)

type MessageType int

type Message struct {
	Seq     int                    `json:"seq"`
	Type    MessageType            `json:"type"`
	From    string                 `json:"from"`
	To      string                 `json:"to"`
	Content map[string]interface{} `json:"content"`
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
