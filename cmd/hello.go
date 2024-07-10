package cmd

import (
	"github.com/gkwa/hollowbeak/core"
	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.V(1).Info("Debug: Entering hello command Run function")
		logger.Info("Running hello command")
		core.Hello(logger)
		logger.V(1).Info("Debug: Exiting hello command Run function")
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	logger.V(1).Info("Debug: Added hello command to root command")
}
