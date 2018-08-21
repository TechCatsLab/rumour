/*
 * Revision History:
 *     Initial: 2018/08/09        Tong Yuehong
 */

package response

import (
	"github.com/TechCatsLab/apix/http/server"
	"github.com/TechCatsLab/rumour/constants"
)

// WriteStatusAndDataJSON send JSON encoding status (and data) response.
// The data is optional.
func WriteStatusAndDataJSON(ctx *server.Context, status int, data interface{}) error {
	if data == nil {
		return ctx.ServeJSON(map[string]interface{}{constants.RespKeyStatus: status})
	}

	return ctx.ServeJSON(map[string]interface{}{
		constants.RespKeyStatus: status,
		constants.RespKeyData:   data,
	})
}

// WriteStatusAndIDJSON send JSON encoding status and id response.
func WriteStatusAndIDJSON(ctx *server.Context, status int, id interface{}) error {
	return ctx.ServeJSON(map[string]interface{}{
		constants.RespKeyStatus: status,
		constants.RespKeyID:     id,
	})
}

