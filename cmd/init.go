// Copyright (c) 2023 John Dewey

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

// Package cmd implements the CLI commands for gilt.
package cmd

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/retr0h/gilt/v2/pkg/config"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Gilt with a Giltfile",
	Long: `Initializes Gilt by creating a default config file in the shell's
current working directory.`,
	Run: func(_ *cobra.Command, _ []string) {
		var b bytes.Buffer

		// set configFile defaults
		repo := []config.Repository{
			{
				Git:     "",
				Version: "",
				DstDir:  "",
			},
		}
		viper.SetDefault("repositories", repo)
		c := viper.AllSettings()

		ye := yaml.NewEncoder(&b)
		ye.SetIndent(2)
		if err := ye.Encode(c); err != nil {
			logFatal(
				"failed to encode file",
				slog.Group("",
					slog.String("err", err.Error()),
				),
			)
		}

		configFile := viper.GetString("giltFile")
		_, err := os.Stat(configFile)
		if err == nil {
			logFatal(
				"file already exists",
				slog.Group("",
					slog.String("Giltfile", configFile),
				),
			)
		}

		if err := os.WriteFile(configFile, b.Bytes(), 0o644); err != nil {
			logFatal(
				"failed to write file",
				slog.Group("",
					slog.String("Giltfile", viper.ConfigFileUsed()),
					slog.String("err", err.Error()),
				),
			)
		}

		fmt.Printf("wrote %s\n", configFile)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
