/*
 * Revision History:
 *     Initial: 2018/07/08        Tong Yuehong
 */

package message

import (
	json "github.com/json-iterator/go"
)

// TextMessage -
type TextMessage struct {
}

// Marshal -
func (text *TextMessage) Marshal() ([]byte, error) {
	return json.Marshal(text)
}

// Unmarshal -
func (text *TextMessage) Unmarshal(info []byte) error {
	return json.Unmarshal(info, &text)
}
