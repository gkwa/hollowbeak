package cmd

import (
	"fmt"
	"os"

	"github.com/gkwa/hollowbeak/core"
	"github.com/spf13/cobra"
)

var outputFormat string

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.V(1).Info("Debug: Entering hello command Run function")
		logger.Info("Running hello command")
		if err := core.Hello(logger, outputFormat); err != nil {
			logger.Error(err, "Failed to execute Hello function")
			os.Exit(1)
		}
		logger.V(1).Info("Debug: Exiting hello command Run function")
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringVar(&outputFormat, "output", "html", "Output format: 'markdown' or 'html'")
	if err := helloCmd.MarkFlagRequired("output"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'output' flag as required: %v\n", err)
		os.Exit(1)
	}
}
