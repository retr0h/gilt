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

package repositories

import (
	"github.com/retr0h/go-gilt/internal/repository"
)

// Repositories representing the Giltfile.yaml file.
type Repositories struct {
	// Debug enable or disable debug option set from CLI.
	Debug bool `mapstruture:"debug"`
	// GiltFile path to Gilt's config file option set from CLI.
	GiltFile string `mapstructure:"giltFile"`
	// GiltDir path to Gilt's clone dir option set from CLI.
	GiltDir string `mapstructure:"giltDir"`
	// Repositories a slice of repository configurations to overlay.
	Repositories []repository.Repository `mapstruture:"repositories"`
}
