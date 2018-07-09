/*
 * Revision History:
 *     Initial: 2018/05/22        Tong Yuehong
 */

package conf

import (
	"fmt"

	"github.com/TechCatsLab/rumour/pkg/log"
	"github.com/spf13/viper"
)

// Conf -
var Conf *Config

// Config is used to config the access endpoint
type Config struct {
	ServerAddr        string
	CorsHosts         []string
	SecretKey         string
	WSReadBufferSize  int
	WSWriteBufferSize int
}

func loadConfiguration() *Config {
	viper.AddConfigPath("../conf")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Error(fmt.Sprintf("Read configuration file with error: %v", err))
		panic(err)
	}

	Conf = &Config{
		ServerAddr:        viper.GetString("serverAddr"),
		CorsHosts:         viper.GetStringSlice("middleware.cors.hosts"),
		SecretKey:         viper.GetString("middleware.secretKey"),
		WSReadBufferSize:  viper.GetInt("websocket.readBufferSize"),
		WSWriteBufferSize: viper.GetInt("websocket.writeBufferSize"),
	}

	return Conf
}

func init() {
	loadConfiguration()
}
