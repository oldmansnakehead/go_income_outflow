// go run main.go db seed

package cmd

import (
	"go_income_outflow/db/migrations/seeds"
	"log"

	"github.com/spf13/cobra"
)

var dbSeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		// Seed
		if err := seeds.UserSeeder(dbConn); err != nil {
			log.Fatal(err.Error())
		}

		if err := seeds.AccountSeeder(dbConn); err != nil {
			log.Fatal(err.Error())
		}

		if err := seeds.TransactionCategorySeeder(dbConn); err != nil {
			log.Fatal(err.Error())
		}

		if err := seeds.CreditCardSeeder(dbConn); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	dbCmd.AddCommand(dbSeedCmd)
}
