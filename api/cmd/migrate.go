/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Kudryavkaz/sztuea-api/internal/log"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database"
	"github.com/Kudryavkaz/sztuea-api/internal/resource/database/model"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "for generating gorm schema migration statements",
	RunE: func(cmd *cobra.Command, args []string) error {
		database.InitDatabase()

		err := model.InitModels()
		if err != nil {
			log.Logger().Error("Models migrate fail.", zap.Error(err))
			return err
		}
		log.Logger().Info("Models migrate success.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
