package cmd

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migrations up or down",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbUser := viper.GetString("DB_USER")
		dbPassword := viper.GetString("DB_PASSWORD")
		dbHost := viper.GetString("DB_HOST")
		dbPort := viper.GetString("DB_PORT")
		database := viper.GetString("DATABASE")
		migrationsDir := viper.GetString("MIGRATIONS_DIR")

		m, err := migrate.New(
			fmt.Sprintf("file://%s", migrationsDir),
			fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, database))

		if err != nil {
			return err
		}

		down, err := cmd.Flags().GetBool("down")
		if err != nil {
			return err
		}

		if down {
			log.Print("Migrating DOWN")
			return m.Down()
		}

		log.Print("Migrating UP")
		return m.Up()
	},
}

func init() {
	dbCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().BoolP("down", "d", false, "migrate down")
}
