package cmds

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/cli/server/api/server"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start communication server",
	Run: func(cmd *cobra.Command, args []string) {
		c := &server.Config{
			User: user,
			Pass: pass,
			Host: host,
			Port: port,
		}

		db, err := c.OpenDB()
		if err != nil {
			log.Error(err)
			return
		}


		err = server.OpenStore(db)
		if err != nil {
			return
		}

		s := server.CreateWebsocketServer()
		if err := s.Serve(); err != nil {
			log.Error("[cmds] Start:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	cobra.OnInitialize(mysqlConfig)

	startCmd.Flags().StringVar(&cfgFile, "config", "", "mysql config file")
}

func mysqlConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AddConfigPath(cfgFile)
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Error("Can't read config:", err)
		os.Exit(1)
	}

	user = viper.GetString("mysql.user")
	pass = viper.GetString("mysql.pass")
	host = viper.GetString("mysql.host")
	port = viper.GetString("mysql.port")
}
