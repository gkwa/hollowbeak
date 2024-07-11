package cmd

import (
	"fmt"
	"os"

	"github.com/gkwa/hollowbeak/core"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	inputFile    string
	fetcherTypes []string
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your command",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.V(1).Info("Debug: Entering hello command Run function")
		logger.Info("Running hello command")
		if err := core.Hello(logger, inputFile, outputFormat, fetcherTypes); err != nil {
			logger.Error(err, "Failed to execute Hello function")
			os.Exit(1)
		}
		logger.V(1).Info("Debug: Exiting hello command Run function")
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringVar(&outputFormat, "output", "markdown", "Output format: 'markdown' or 'html'")
	helloCmd.Flags().StringVar(&inputFile, "input", "", "Input file path")
	helloCmd.Flags().StringSliceVar(&fetcherTypes, "fetcher", []string{"colly"}, "Title fetcher types: 'http', 'colly', or 'sql'. Can be specified multiple times.")
	if err := helloCmd.MarkFlagRequired("input"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking 'input' flag as required: %v\n", err)
		os.Exit(1)
	}
}
