package group

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const inexistentFilePath string = "/foo/bar/yadayada"
const gshadowFileMockPath string = "gshadow_mock.txt"

func TestReadShadowDB(t *testing.T) {
	groups, err := readShadowDB(gshadowFileMockPath)
	assert.Nil(t, err, "No error should be returned")
	assert.Equal(t, 1, len(*groups), "Only one group has a real password")
	_, ok := (*groups)["foobar"]
	assert.True(t, ok, "Group has the expected name")
}

func TestReadShadowDBInexistentFile(t *testing.T) {
	groups, err := readShadowDB(inexistentFilePath)
	assert.NotNil(t, err, "Must get an error")
	assert.Contains(t, err.Error(), "no such file")
	assert.Equal(t, 0, len(*groups), "The returned map is empty")
}

func TestReadDBFromFileInvalid(t *testing.T) {
	groups, err := ReadDBFromFile(100, 100, inexistentFilePath, "/etc/shadow")
	assert.NotNil(t, err, "Must get an error")
	assert.Contains(t, err.Error(), "no such file")
	assert.Equal(t, 0, len(groups), "The returned map is empty")
}

func TestReadDBFromFile(t *testing.T) {
	groups, err := ReadDBFromFile(130, 2000, "group_mock.txt", gshadowFileMockPath)
	assert.Nil(t, err, "No error should be returned")
	assert.Equal(t, 8, len(groups), "The expected number of groups was returned")
	last := len(groups) - 1
	assert.Equal(t, "foobar", groups[last].name, "Got the expected name for last group")

	for _, group := range groups[:last] {
		assert.Equal(t, "", group.password, "Almost all groups have no password")
	}

	assert.NotEqual(t, "", groups[last].password, "Only the last group has a password")
}
