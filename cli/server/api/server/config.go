/*
 * Revision History:
 *     Initial: 2018/08/22        Tong Yuehong
 */

package server

import (
	"fmt"
	"database/sql"

	log "github.com/TechCatsLab/logging/logrus"
)

type Config struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

func (c *Config) OpenDB() (*sql.DB, error) {
	dataSource := fmt.Sprintf(c.User + ":" + c.Pass + "@" + "tcp(" + c.Host + ":" + c.Port + ")/" + "chat" + "?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return db, nil
}
