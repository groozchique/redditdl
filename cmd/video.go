package cmd

import (
	"github.com/handsomefox/redditdl/pkg/downloader/config"
	"github.com/spf13/cobra"
)

var videoCmd = &cobra.Command{
	Use:   "video",
	Short: "Command for downloading videos",
	Long: `This command allows users to download video files from
different subreddits from reddit.com, applying different
filters to the content which will be downloaded.
`,
	Run: func(cmd *cobra.Command, _ []string) {
		settings := GetSettings(cmd)
		settings.ContentType = config.ContentVideos
		RunCommand(&settings)
	},
}

func init() {
	rootCmd.AddCommand(videoCmd)
	SetCommonFlags(videoCmd)
}
