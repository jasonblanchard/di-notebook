package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "interact with db",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("Use a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.PersistentFlags().String("migrationsDir", os.Getenv("MIGRATIONS_DIR"), "migrations directory")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
