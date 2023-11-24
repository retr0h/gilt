// Copyright (c) 2018 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

package cmd

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/retr0h/go-gilt/internal"
	"github.com/retr0h/go-gilt/internal/config"
	"github.com/retr0h/go-gilt/internal/exec"
	"github.com/retr0h/go-gilt/internal/git"
	"github.com/retr0h/go-gilt/internal/repositories"
	"github.com/retr0h/go-gilt/internal/repository"
)

var (
	repos     internal.RepositoriesManager
	logger    *slog.Logger
	appConfig config.Repositories
	appFs     afero.Fs
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-gilt",
	Short: "A GIT layering command line tool",
	Long: `
           o  o
         o |  |
    o--o   | -o-
    |  | | |  |
    o--O | o  o
       |
    o--o

A GIT layering command line tool.

https://github.com/retr0h/go-gilt
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger) // initial logger
	cobra.OnInitialize(initConfig, initLogger)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable or disable debug mode")
	rootCmd.PersistentFlags().
		StringP("gilt-dir", "c", "~/.gilt/clone", "Path to Gilt's clone dir")
	rootCmd.PersistentFlags().
		StringP("gilt-file", "f", "Giltfile.yaml", "Path to config file")

	_ = viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	_ = viper.BindPFlag("giltFile", rootCmd.PersistentFlags().Lookup("gilt-file"))
	_ = viper.BindPFlag("giltDir", rootCmd.PersistentFlags().Lookup("gilt-dir"))
	_ = viper.BindPFlag("repositories", rootCmd.PersistentFlags().Lookup("repositories"))
}

func initLogger() {
	logLevel := slog.LevelInfo
	if viper.GetBool("debug") {
		logLevel = slog.LevelDebug
	}

	logger = slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      logLevel,
			TimeFormat: time.Kitchen,
		}),
	)
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("gilt")
	viper.SetConfigFile(viper.GetString("giltFile"))

	if err := viper.ReadInConfig(); err != nil {
		logger.Error(
			"failed to read config",
			slog.String("Giltfile", viper.ConfigFileUsed()),
			slog.String("err", err.Error()),
		)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&appConfig); err != nil {
		logger.Error(
			"failed to unmarshal config",
			slog.String("Giltfile", viper.ConfigFileUsed()),
			slog.String("err", err.Error()),
		)
		os.Exit(1)
	}

	appFs = afero.NewOsFs()

	repoManager := repository.NewCopy(
		appFs,
		logger,
	)

	execManager := exec.New(
		appConfig.Debug,
		logger,
	)

	g := git.New(
		appFs,
		appConfig.Debug,
		execManager,
		logger,
	)

	repos = repositories.New(
		appFs,
		appConfig,
		repository.New(
			appFs,
			repoManager,
			g,
			logger,
		),
		execManager,
		logger,
	)
}
