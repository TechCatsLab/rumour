/*
 * Revision History:
 *     Initial: 2018/05/30        Tong Yuehong
 */

package authenticate

import (
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/TechCatsLab/rumour/conf"
)

// Authenticate - authenticate user's token.
type Authenticate struct {
}

// Authenticate -
func (auth *Authenticate) Authenticate(t string) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Conf.SecretKey), nil
	}

	token, err := jwt.Parse(t, keyFunc)
	if err == nil && token.Valid {
		return token.Claims.(jwt.MapClaims)["uid"].(string), nil
	}
	return "", err
}
