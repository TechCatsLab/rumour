/*
 * Revision History:
 *     Initial: 2018/08/09        Tong Yuehong
 */

package constants

const (
	// ErrSucceed - Succeed
	ErrSucceed = 0

	// ErrPermission - Permission Denied
	ErrPermission = 401
	ErrForbidden  = 438

	// ErrToken - Invalid Token
	ErrToken = 420

	// ErrInvalidParam - Invalid Parameter
	ErrInvalidParam = 421

	// ErrAccount - No This User or Password Error
	ErrAccount = 422

	// ErrSubNats - Subscribe to nats error
	ErrSubNats = 423

	// ErrInternalServerError - Internal error.
	ErrInternalServerError = 500

	// ErrWechatPay - Wechat Pay error.
	ErrWechatPay = 520

	// ErrWechatAuth - Wechat Auth error.
	ErrWechatAuth = 521

	// ErrMongoDB - MongoDB operations error.
	ErrMongoDB = 600

	// ErrMysql - Mysql operations error.
	ErrMysql = 700

	ErrDuplicate = 800
	ErrNotFound  = 801
)
