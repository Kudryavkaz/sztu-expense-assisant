package cmd

import (
	"github.com/Kudryavkaz/sztuea-api/internal/api"
	"github.com/samber/lo"
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
		lo.Must0(api.StartSever(port, prefork))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().Uint16VarP(&port, "port", "p", 3000, "Server Port")
	startCmd.Flags().BoolVar(&prefork, "prefork", true, "Use of the SO_REUSEPORT socket option")
}
