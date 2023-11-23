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

// TODO(retr0h): move to public testing pattern
package path

import (
	"fmt"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PathTestSuite struct {
	suite.Suite
}

func (suite *PathTestSuite) TestexpandUserOk() {
	originalCurrentUser := currentUser
	currentUser = func() (*user.User, error) {
		return &user.User{
			HomeDir: "/testUser",
		}, nil
	}
	defer func() { currentUser = originalCurrentUser }()

	got, err := ExpandUser("~/foo/bar")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), got, "/testUser/foo/bar")
}

func (suite *PathTestSuite) TestexpandUserReturnsError() {
	originalCurrentUser := currentUser
	currentUser = func() (*user.User, error) {
		return nil, fmt.Errorf("failed to get current user")
	}
	defer func() { currentUser = originalCurrentUser }()

	_, err := ExpandUser("~/foo/bar")
	assert.Error(suite.T(), err)
}

// In order for `go test` to run this suite, we need to create
// a normal test function and pass our suite to suite.Run.
func TestPathTestSuite(t *testing.T) {
	suite.Run(t, new(PathTestSuite))
}
