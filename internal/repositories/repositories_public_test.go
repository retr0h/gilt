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

package repositories_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/retr0h/go-gilt/internal/git"
	"github.com/retr0h/go-gilt/internal/repositories"
	helper "github.com/retr0h/go-gilt/internal/testing"
)

type RepositoriesTestSuite struct {
	suite.Suite
	r repositories.Repositories
}

func (suite *RepositoriesTestSuite) unmarshalYAML(data []byte) error {
	return helper.UnmarshalYAML([]byte(data), &suite.r.Repositories)
}

func (suite *RepositoriesTestSuite) SetupTest() {
	suite.r = repositories.Repositories{}
	suite.r.GiltDir = helper.CreateTempDirectory()
}

func (suite *RepositoriesTestSuite) TearDownTest() {
	helper.RemoveTempDirectory(suite.r.GiltDir)
}

func (suite *RepositoriesTestSuite) TestOverlayFailsCloneReturnsError() {
	data := `
---
- git: invalid.
  version: abc1234
  dstDir: path/user.repo
`
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	anon := func() error {
		err := suite.r.Overlay()
		assert.Error(suite.T(), err)

		return err
	}

	git.MockRunCommandErrorsOn("git", anon)
}

func (suite *RepositoriesTestSuite) TestOverlayFailsCheckoutIndexReturnsError() {
	data := `
---
- git: https://example.com/user/repo.git
  version: abc1234
  dstDir: /invalid/directory
`
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	anon := func() error {
		err := suite.r.Overlay()
		assert.Error(suite.T(), err)

		return err
	}

	git.MockRunCommandErrorsOn("checkout-index", anon)
}

func (suite *RepositoriesTestSuite) TestOverlay() {
	data := `
---
- git: https://example.com/user/repo1.git
  version: abc1234
  dstDir: path/user.repo
- git: https://example.com/user/repo2.git
  version: abc1234
  sources:
    - src: foo
      dstFile: bar
`
	err := suite.unmarshalYAML([]byte(data))
	assert.NoError(suite.T(), err)

	anon := func() error {
		err := suite.r.Overlay()
		assert.NoError(suite.T(), err)

		return err
	}

	dstDir, _ := git.FilePathAbs(suite.r.Repositories[0].DstDir)
	got := git.MockRunCommand(anon)
	want := []string{
		fmt.Sprintf(
			"git clone https://example.com/user/repo1.git %s/https---example.com-user-repo1.git-abc1234",
			suite.r.GiltDir,
		),
		fmt.Sprintf("git -C %s/https---example.com-user-repo1.git-abc1234 reset --hard abc1234",
			suite.r.GiltDir),
		fmt.Sprintf(
			"git -C %s/https---example.com-user-repo1.git-abc1234 checkout-index --force --all --prefix %s",
			suite.r.GiltDir,
			(dstDir + string(os.PathSeparator)),
		),
		fmt.Sprintf(
			"git clone https://example.com/user/repo2.git %s/https---example.com-user-repo2.git-abc1234",
			suite.r.GiltDir,
		),
		fmt.Sprintf("git -C %s/https---example.com-user-repo2.git-abc1234 reset --hard abc1234",
			suite.r.GiltDir),
	}

	assert.Equal(suite.T(), want, got)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestRepositoriesTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesTestSuite))
}
