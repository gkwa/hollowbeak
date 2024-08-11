package cmd

import (
	"bytes"
	"os"
	"strings"

	"github.com/gkwa/hollowbeak/core"
	"github.com/spf13/cobra"
)

// fetchUrlNamesCmd represents the fetchUrlNames command
var fetchUrlNamesCmd = &cobra.Command{
	Use:     "fetch-url-names url [url...]",
	Short:   "Fetch url names given a list of urls",
	Aliases: []string{"fun"},
	Args:    cobra.MinimumNArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		logger.V(1).Info("Debug: Entering hello command Run function")
		logger.Info("Running hello command")

		buffer := new(bytes.Buffer)
		content := strings.Join(args, "\n")
		buffer.WriteString(content)

		if err := core.FetchURLTitles(
			logger,
			buffer,
			outputFormat,
			fetcherTypes,
			noCache,
		); err != nil {
			logger.Error(err, "Failed to execute Hello function")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchUrlNamesCmd)
}
