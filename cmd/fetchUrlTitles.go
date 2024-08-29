package cmd

import (
	"os"

	"github.com/gkwa/hollowbeak/core"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	fetcherTypes []string
	noCache      bool
)

var fileUrlTitlesCmd = &cobra.Command{
	Use:     "file-url-titles",
	Short:   "Accept path to file and return urls and url titles in various formats (default markdown)",
	Aliases: []string{"efu"},
	Args:    cobra.ExactArgs(1),
	Long:    `A longer description that spans multiple lines and likely contains examples and usage of using your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.V(1).Info("Debug: Entering hello command Run function")
		logger.Info("Running hello command")

		file, err := os.Open(args[0])
		if err != nil {
			logger.Error(err, "Failed to open input file")
			os.Exit(1)
		}
		defer file.Close()

		if err := core.FetchURLTitles(
			logger,
			file,
			outputFormat,
			fetcherTypes,
			noCache,
		); err != nil {
			logger.Error(err, "Failed to execute Hello function")
			os.Exit(1)
		}
		logger.V(1).Info("Debug: Exiting hello command Run function")
	},
}

func init() {
	rootCmd.AddCommand(fileUrlTitlesCmd)
	fileUrlTitlesCmd.Flags().StringSliceVar(&fetcherTypes, "fetcher", []string{"sql", "colly", "http"}, "Title fetcher types: 'http', 'colly', or 'sql'. Can be specified multiple times.")
	fileUrlTitlesCmd.Flags().BoolVar(&noCache, "no-cache", false, "Skip cache for this run")
}
