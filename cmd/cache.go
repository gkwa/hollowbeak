package cmd

import (
	"fmt"

	"github.com/gkwa/hollowbeak/core"
	"github.com/spf13/cobra"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Display cacheuration information",
	Run: func(cmd *cobra.Command, args []string) {
		displaycachePath()
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}

func displaycachePath() {
	cachePath, err := core.GetCachePath()
	if err != nil {
		fmt.Printf("Error getting cache path: %v\n", err)
		return
	}
	fmt.Println(cachePath)
}
