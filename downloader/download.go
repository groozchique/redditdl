package downloader

import (
	"time"

	"github.com/handsomefox/redditdl/internal/logging"
	"github.com/handsomefox/redditdl/internal/utils"
)

type MediaType uint8

const (
	_ MediaType = iota
	MediaImages
	MediaVideos
	MediaAny
)

// Settings is the configuration for the Downloader.
type Settings struct {
	Directory   string
	Subreddit   string
	Sorting     string
	Timeframe   string
	Orientation string

	Count     int64
	MinWidth  int
	MinHeight int

	Verbose      bool
	ShowProgress bool

	ContentType MediaType
}

const (
	SleepTime   = 5 * time.Second // sleep period between post fetches
	workerCount = 16              // amount of goroutines for downloading files.
)

// Download downloads the images according to the given configuration.
func Download(settings *Settings, filters []Filter) (int64, error) {
	dl := downloader{
		client:  utils.CreateClient(),
		log:     logging.GetLogger(settings.Verbose),
		counter: counter{},
	}

	return dl.download(settings, filters)
}
