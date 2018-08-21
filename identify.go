/*
 * Revision History:
 *     Initial: 2018/05/30        Tong Yuehong
 */

package rumour

import (
	"errors"
)

var ErrInvalidID = errors.New("invalid identify")

// Identify - represent user's identify.
type Identify string

// Identify - return userid.
func (identify Identify) Id() (string, error) {
	if string(identify) == "" {
		return "", ErrInvalidID
	}
	return string(identify), nil
}

// Equal
func (identify Identify) Equal(id Identify) bool {
	d, _ := id.Id()

	return identify == Identify(d)
}
