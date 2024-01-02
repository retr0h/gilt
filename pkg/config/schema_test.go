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

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/suite"
)

type SchemaTestSuite struct {
	suite.Suite

	v *validator.Validate
}

func (suite *SchemaTestSuite) SetupTest() {
	suite.v = validator.New()
	err := registerValidators(suite.v)
	assert.NoError(suite.T(), err)
}

func (suite *SchemaTestSuite) TestRepositories() {
	tests := []struct {
		param    *Repositories
		expected string
	}{
		{&Repositories{
			GiltFile: "giltFile",
			GiltDir:  "giltDir",
			Repositories: []Repository{
				{
					Git:     "gitURL",
					Version: "abc1234",
					DstDir:  "dstDir",
				},
			},
		}, ""},
		{&Repositories{
			GiltFile: "giltFile",
			GiltDir:  "",
			Repositories: []Repository{
				{
					Git:     "gitURL",
					Version: "abc1234",
					DstDir:  "dstDir",
				},
			},
		}, "Key: 'Repositories.GiltDir' Error:Field validation for 'GiltDir' failed on the 'required' tag"},
		{&Repositories{
			GiltFile: "",
			GiltDir:  "giltDir",
			Repositories: []Repository{
				{
					Git:     "gitURL",
					Version: "abc1234",
					DstDir:  "dstDir",
				},
			},
		}, "Key: 'Repositories.GiltFile' Error:Field validation for 'GiltFile' failed on the 'required' tag"},
		{&Repositories{
			GiltFile: "giltFile",
			GiltDir:  "giltDir",
		}, "Key: 'Repositories.Repositories' Error:Field validation for 'Repositories' failed on the 'required' tag"},
	}

	// NOTE(nic): we have an entrypoint for validating this schema, so use it to
	// ensure test coverage.  All the other tests will use the validator in the suite.
	for _, test := range tests {
		err := Validate(test.param)
		if test.expected != "" {
			assert.EqualError(suite.T(), err, test.expected)
		} else {
			assert.NoError(suite.T(), err)
		}
	}
}

func (suite *SchemaTestSuite) TestSourceSchema() {
	tests := []struct {
		param    *Source
		expected string
	}{
		{&Source{
			Src:     "src",
			DstFile: "dstFile",
		}, ""},
		{&Source{
			Src:    "src",
			DstDir: "dstDir",
		}, ""},
		{&Source{
			Src: "",
		}, "Key: 'Source.Src' Error:Field validation for 'Src' failed on the 'required' tag\nKey: 'Source.DstFile' Error:Field validation for 'DstFile' failed on the 'required_without' tag\nKey: 'Source.DstDir' Error:Field validation for 'DstDir' failed on the 'required_without' tag"},
		{&Source{
			Src:     "src",
			DstFile: "dstFile",
			DstDir:  "dstDir",
		}, "Key: 'Source.DstFile' Error:Field validation for 'DstFile' failed on the 'excluded_with' tag\nKey: 'Source.DstDir' Error:Field validation for 'DstDir' failed on the 'excluded_with' tag"},
	}

	for _, test := range tests {
		err := suite.v.Struct(test.param)
		if test.expected != "" {
			assert.EqualError(suite.T(), err, test.expected)
		} else {
			assert.NoError(suite.T(), err)
		}
	}
}

func (suite *SchemaTestSuite) TestCommandSchema() {
	tests := []struct {
		param    *Command
		expected string
	}{
		{&Command{
			Cmd: "foo",
		}, ""},
		{&Command{
			Cmd:  "foo",
			Args: []string{"bar", "baz"},
		}, ""},
		{&Command{
			Cmd: "",
		}, "Key: 'Command.Cmd' Error:Field validation for 'Cmd' failed on the 'required' tag"},
		{&Command{
			Args: []string{"bar", "baz"},
		}, "Key: 'Command.Cmd' Error:Field validation for 'Cmd' failed on the 'required' tag"},
	}

	for _, test := range tests {
		err := suite.v.Struct(test.param)

		if test.expected != "" {
			assert.EqualError(suite.T(), err, test.expected)
		} else {
			assert.NoError(suite.T(), err)
		}
	}
}

func (suite *SchemaTestSuite) TestRepositorySchema() {
	tests := []struct {
		param    *Repository
		expected string
	}{
		{&Repository{
			Git:     "gitURL",
			Version: "abc1234",
			DstDir:  "dstDir",
		}, ""},
		{&Repository{
			Git:     "gitURL",
			Version: "v1.1",
			DstDir:  "dstDir",
		}, ""},
		{&Repository{
			Git:     "gitURL",
			Version: "abc1234",
			Sources: []Source{
				{
					Src:     "src",
					DstFile: "dstFile",
				},
			},
		}, ""},
		{&Repository{
			Git:     "gitURL",
			Version: "",
			DstDir:  "dstDir",
		}, "Key: 'Repository.Version' Error:Field validation for 'Version' failed on the 'required' tag"},
		{&Repository{
			Git:     "gitURL",
			Version: "abc1234",
			DstDir:  "dstDir",
			Sources: []Source{
				{
					Src:     "src",
					DstFile: "dstFile",
				},
			},
		}, "Key: 'Repository.DstDir' Error:Field validation for 'DstDir' failed on the 'excluded_with' tag\nKey: 'Repository.Sources[0]' Error:Field validation for 'Sources[0]' failed on the 'excluded_with' tag"},
	}

	for _, test := range tests {
		err := suite.v.Struct(test.param)

		if test.expected != "" {
			assert.EqualError(suite.T(), err, test.expected)
		} else {
			assert.NoError(suite.T(), err)
		}
	}
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestSchemaTestSuite(t *testing.T) {
	suite.Run(t, new(SchemaTestSuite))
}
