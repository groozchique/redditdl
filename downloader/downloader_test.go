package downloader_test

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/handsomefox/redditdl/cmd"
	"github.com/handsomefox/redditdl/cmd/params"
	"github.com/handsomefox/redditdl/downloader"
)

func TestDownload(t *testing.T) {
	t.Parallel()

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	p := setupConfig(dir, 25)

	dl, err := downloader.New(p, downloader.DefaultFilters()...)
	if err != nil {
		t.Fatal(err)
	}

	status := dl.Download(context.TODO())
	if status.Finished < 1 {
		t.Error("Failed to download requested amount", status.Finished, p.MediaCount)
	}
}

func setupConfig(dir string, count int64) *params.CLIParameters {
	cmd.SetGlobalLoggingLevel(false)
	cliParams := &params.CLIParameters{
		Sort:             "best",
		Timeframe:        "all",
		Directory:        dir,
		Subreddits:       []string{"wallpaper"},
		MediaMinWidth:    0,
		MediaMinHeight:   0,
		MediaCount:       count,
		MediaOrientation: params.RequiredOrientationAny,
		ContentType:      params.RequiredContentTypeImages,
		ShowProgress:     false,
		VerboseLogging:   false,
	}
	return cliParams
}

func BenchmarkDownload10(b *testing.B) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	p := setupConfig(dir, 10)
	dl, err := downloader.New(p, downloader.DefaultFilters()...)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_ = dl.Download(context.TODO())
	}
}

func BenchmarkDownload50(b *testing.B) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	p := setupConfig(dir, 50)
	dl, err := downloader.New(p, downloader.DefaultFilters()...)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_ = dl.Download(context.TODO())
	}
}

func TestNewFilename(t *testing.T) {
	t.Parallel()
	type args struct {
		name      string
		extension string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Create file.jpg",
			args: args{
				name:      "file",
				extension: "jpg",
			},
			want:    "file.jpg",
			wantErr: false,
		}, {
			name: "Create a file with invalid characters",
			args: args{
				name:      "/<>:file",
				extension: "/<>:jpg",
			},
			want:    "file.jpg",
			wantErr: false,
		}, {
			name: "Empty name",
			args: args{
				name:      "",
				extension: "jpg",
			},
			want:    "",
			wantErr: true,
		}, {
			name: "Empty extension",
			args: args{
				name:      "file",
				extension: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := downloader.NewFormattedFilename(tt.args.name, tt.args.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	t.Parallel()
	exec, err := os.Executable()
	if err != nil {
		t.Fatalf("couldn't find the running executable")
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check existing file",
			args: args{
				filename: exec,
			},
			want: true,
		}, {
			name: "Check non-existing file",
			args: args{
				filename: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := downloader.FileExists(tt.args.filename); got != tt.want {
				t.Errorf("Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNavigateTo(t *testing.T) {
	t.Parallel()
	type args struct {
		dir       string
		createDir bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Navigate to created directory",
			args: args{
				dir:       path.Join(os.TempDir(), "test_dir"),
				createDir: true,
			},
			wantErr: false,
		}, {
			name: "Navigate to non-existing directory",
			args: args{
				dir:       "/<>:",
				createDir: false,
			},
			wantErr: true,
		}, {
			name: "Navigate to existing directory",
			args: args{
				dir:       os.TempDir(),
				createDir: false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := downloader.ChdirOrCreate(tt.args.dir, tt.args.createDir); (err != nil) != tt.wantErr {
				t.Errorf("NavigateTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	t.Parallel()
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "http://",
			args: args{
				str: "http://",
			},
			want: false,
		},
		{
			name: "google.com",
			args: args{
				str: "google.com",
			},
			want: false,
		},
		{
			name: "google",
			args: args{
				str: "google",
			},
			want: false,
		},
		{
			name: "www.google",
			args: args{
				str: "www.google",
			},
			want: false,
		},
		{
			name: "http://google.com",
			args: args{
				str: "http://google.com",
			},
			want: true,
		},
		{
			name: "https://google.com",
			args: args{
				str: "https://google.com",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := downloader.IsValidURL(tt.args.str); got != tt.want {
				t.Errorf("IsURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
