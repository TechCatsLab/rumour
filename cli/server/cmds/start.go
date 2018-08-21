package cmds

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/cobra"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/cli/server/api/server"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start communication server",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.CreateWebsocketServer()

		if err := s.Serve(); err != nil {
			log.Error("[cmds] Start:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
