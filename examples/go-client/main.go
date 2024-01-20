package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"

	"github.com/retr0h/gilt/pkg/config"
	"github.com/retr0h/gilt/pkg/repositories"
)

type repositoriesManager interface {
	Overlay() error
}

func getLogger(debug bool) *slog.Logger {
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.Kitchen,
		}),
	)

	return logger
}

func main() {
	debug := true
	logger := getLogger(debug)

	c := config.Repositories{
		Debug:   debug,
		GiltDir: "~/.gilt",
		Repositories: []config.Repository{
			{
				Git:     "https://github.com/retr0h/ansible-etcd.git",
				Version: "77a95b7",
				DstDir:  "../../test/integration/tmp/retr0h.ansible-etcd",
			},
		},
	}

	var r repositoriesManager = repositories.New(c, logger)
	r.Overlay()
}
