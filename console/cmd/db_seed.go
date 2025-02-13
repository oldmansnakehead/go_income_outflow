package cmd

import (
	"go_income_outflow/db"
	"go_income_outflow/db/migrations/seeds"
	"log"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

var dbSeedCmd = &cobra.Command{
	Use:   "dbSeed",
	Short: "Seed the database with initial data",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		dbConn = db.LazyConnect()
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Seed
		if err := seeds.UserSeeder(dbConn); err != nil {
			log.Fatal(err.Error())
		}

		if err := seeds.AccountSeeder(dbConn); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(dbSeedCmd)
}
