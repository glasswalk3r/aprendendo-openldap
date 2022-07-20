package passwd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const mockDBFilename string = "mock.txt"

func TestReadDBFromFile(t *testing.T) {
	users, err := ReadDBFromFile(1000, 2000, 1000, 2000, mockDBFilename)
	assert.Nilf(t, err, "Unexpected error: '%s'", err)
	assert.Equalf(t, 1, len(users), "Only one user is expected, got: %v", users)
  expected := DBEntry{
    "alceu",
    1000,
    1000,
    GECOS{"Alceu Rodrigues de Freitas Junior","571","+551155422748","+551155422748"},
    "/home/alceu",
    "/bin/bash",
  }
  assert.EqualValues(t, expected, users[0])
}

func TestReadDBFromFileInvalidFilePath(t *testing.T) {
  users, err := ReadDBFromFile(1000, 2000, 1000, 2000, "/foobar/yadayada")
  assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such file")
	assert.Equalf(t, 0, len(users), "Unexpected content: %v", users)
}
