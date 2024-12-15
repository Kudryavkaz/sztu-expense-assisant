package cmd

import (
	"github.com/Kudryavkaz/sztuea-api/internal/api"
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/cache"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database"
	"github.com/spf13/cobra"
)

var (
	port    uint16
	prefork bool
)

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start the HTTP server",
	Example: "sztuea-api start --port 3000",
	RunE: func(cmd *cobra.Command, args []string) error {
		database.InitDatabase()

		cache.InitRedis()

		if err := api.StartSever(port, prefork); err != nil {
			log.Logger().Error("Failed to start the server")
			panic(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().Uint16VarP(&port, "port", "p", 3000, "Server Port")
	startCmd.Flags().BoolVar(&prefork, "prefork", false, "Use of the SO_REUSEPORT socket option")
}
