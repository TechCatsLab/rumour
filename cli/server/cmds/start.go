package cmds

import (
	"github.com/spf13/cobra"

	"github.com/TechCatsLab/rumour/cli/server/api/server"
	"github.com/TechCatsLab/rumour/pkg/log"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start communication server",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.CreateWebsocketServer()

		if err := s.Serve(); err != nil {
			log.Error("[cmds] Start:", log.Err(err))
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
