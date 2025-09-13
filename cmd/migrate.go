package cmd

import (
	"github.com/LocTr/my-blog-be/database/migrations"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run bun migration tool",
	Run: func(cmd *cobra.Command, args []string) {
		migrations.Migrate()
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
}
