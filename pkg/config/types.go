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

package config

// Repositories perform repository operations.
type Repositories struct {
	// Debug enable or disable debug option set from CLI.
	Debug bool `mapstruture:"debug"`
	// Parallel enable or disable concurrent clone fetches.
	Parallel bool `                           mapstructure:"parallel"`
	// SkipCommands run post-commands as part of the overlay process
	SkipCommands bool
	// GiltFile path to Gilt's config file option set from CLI.
	GiltFile string `                           mapstructure:"giltFile" validate:"required"`
	// GiltDir path to Gilt's clone dir option set from CLI.
	GiltDir string `                           mapstructure:"giltDir"  validate:"required"`
	// Repositories a slice of repository configurations to overlay.
	Repositories []Repository `mapstruture:"repositories"                         validate:"required,dive"`
}

// Source mapping of files and/or directories needing copied.
type Source struct {
	// Src source file or directory to copy.
	Src string `mapstructure:"src"     validate:"required"`
	// DstFile destination of file copy.
	DstFile string `mapstructure:"dstFile" validate:"required_without=DstDir,excluded_with=DstDir"`
	// DstDir destination of directory copy.
	DstDir string `mapstructure:"dstDir"  validate:"required_without=DstFile,excluded_with=DstFile,ne=.,ne=.."`
}

//  Water string `validate:"required_without=Fire,excluded_with=Fire"`

// Command command to execute.
type Command struct {
	Cmd  string   `mapstructure:"cmd"  validate:"required"`
	Args []string `mapstructure:"args"`
}

// Repository contains the repository's details for cloning.
type Repository struct {
	// Git url of Git repository to clone.
	Git string `mapstructure:"git"      validate:"required"`
	// Version the commit SHA or tag to use.
	Version string `mapstructure:"version"  validate:"required"`
	// DstDir destination directory to copy clone to.
	DstDir string `mapstructure:"dstDir"   validate:"required_without=Sources,excluded_with=Sources,ne=.,ne=.."`
	// Sources containing files and/or directories to copy.
	Sources []Source `mapstructure:"sources"  validate:"dive,required_without=DstDir,excluded_with=DstDir"`
	// Commands commands to execute on Repository.
	Commands []Command `mapstructure:"commands"`
}
