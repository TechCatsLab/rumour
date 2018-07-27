/**
 * Revision History:
 *     Initial: 2018/05/22        Tong Yuehong
 */

package rumour

// Endpoint -
type Endpoint interface {
	Serve() error
}
