package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var checkPort uint16

var healthCheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check the api service health",
	Run: func(cmd *cobra.Command, args []string) {
		res, err := http.Get(fmt.Sprintf("http://%s:%d%s", "localhost", checkPort, "/ping"))
		if err != nil {
			fmt.Println("health check failed: ", err.Error())
			os.Exit(1)
		}
		if res.StatusCode != http.StatusOK {
			fmt.Println("health check failed: status code is not 200")
			os.Exit(1)
		}
		fmt.Println("health check success")
	},
}

func init() {
	rootCmd.AddCommand(healthCheckCmd)

	healthCheckCmd.Flags().Uint16VarP(&checkPort, "port", "p", 3000, "Server Port")
}
