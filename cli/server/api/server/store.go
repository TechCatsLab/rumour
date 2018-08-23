/*
 * Revision History:
 *     Initial: 2018/08/23        Tong Yuehong
 */

package server

import (
	"database/sql"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
)

func OpenStore(db *sql.DB) error {
	return mysql.NewStore(db)
}

func CreateTable() error {
	err := mysql.StoreService.Store().Channel().Create()
	if err != nil {
		log.Error(err)
		return err
	}

	err = mysql.StoreService.Store().ChannelUser().Create()
	if err != nil {
		log.Error(err)
		return err
	}

	err = mysql.StoreService.Store().SingleMessage().Create()
	if err != nil {
		log.Error(err)
		return err
	}

	err = mysql.StoreService.Store().ChannelMessage().Create()
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
