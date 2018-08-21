/*
 * Revision History:
 *     Initial: 2018/08/14        Tong Yuehong
 */
package cmds

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"

	log "github.com/TechCatsLab/logging/logrus"
	"github.com/TechCatsLab/rumour/pkg/store/mysql"
)

var (
	host string
	port string
	user string
	pass string
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "database init",
	Run: func(cmd *cobra.Command, args []string) {
		dataSource := fmt.Sprintf(user + ":" + pass + "@" + "tcp(" + host + ":" + port + ")/" + "chat" + "?charset=utf8&parseTime=True&loc=Local")
		db, err := sql.Open("mysql", dataSource)
		if err != nil {
			log.Error(err)
		}

		err = mysql.NewStore(db)
		if err != nil {
			log.Error(err)
		}

		err = mysql.StoreService.Store().Channel().Create()
		if err != nil {
			log.Error(err)
		}

		err = mysql.StoreService.Store().ChannelUser().Create()
		if err != nil {
			log.Error(err)
		}

		err = mysql.StoreService.Store().SingleMessage().Create()
		if err != nil {
			log.Error(err)
		}

		err = mysql.StoreService.Store().ChannelMessage().Create()
		if err != nil {
			log.Error(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&host, "host", "t", "", "host address")
	deployCmd.Flags().StringVarP(&port, "port", "p", "", "port")
	deployCmd.Flags().StringVarP(&user, "user", "u", "", "user")
	deployCmd.Flags().StringVarP(&pass, "pass", "s", "", "pass")
}
