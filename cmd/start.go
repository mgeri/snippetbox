package cmd

import (
	"github.com/mgeri/snippetbox/conf"
	"github.com/mgeri/snippetbox/server"
	"github.com/spf13/cobra"
)

// Version command
func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start " + conf.Executable,
		Long:  "Start " + conf.Executable,
		Run: func(cmd *cobra.Command, args []string) {
			server.ListenAndServe(&logger)
		},
	})
}
