/*
 * Revision History:
 *     Initial: 2018/08/14        Tong Yuehong
 */
package cmds

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/cli/server/api/server"
)

var (
	cfgFile string
	user string
	pass string
	host string
	port string
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "database init",
	Run: func(cmd *cobra.Command, args []string) {
		c := &server.Config {
			User : user,
			Pass: pass,
			Host: host,
			Port: port,
		}

		db, err  := c.OpenDB()
		if err != nil {
			log.Error(err)
			return
		}

		err = server.OpenStore(db)
		if err != nil {
			log.Error(err)
			return
		}

		err = server.CreateTable()
		if err != nil {
			log.Error(err)
			return
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVar(&cfgFile, "config", "", "mysql config file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AddConfigPath(cfgFile)
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	user = viper.GetString("mysql.user")
	pass = viper.GetString("mysql.pass")
	host = viper.GetString("mysql.host")
	port = viper.GetString("mysql.port")
}
