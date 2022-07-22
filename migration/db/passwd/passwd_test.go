package passwd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const mockDBFilename string = "mock.txt"

func TestReadDBFromFile(t *testing.T) {
	users, err := ReadDBFromFile(1000, 2000, 1000, 2000, mockDBFilename)
	assert.Nilf(t, err, "Unexpected error: '%s'", err)
	assert.Equalf(t, 1, len(users), "Only one user is expected, got: %v", users)
	expected := DBEntry{
		User: "alceu",
		UID: 1000,
		GID: 1000,
		GECOS: GECOS{
			FullName: "Alceu Rodrigues de Freitas Junior",
			Office: "571",
			WorkPhone: "+551155422748",
			HomePhone: "+551155422748",
		},
		HomeDir: "/home/alceu",
		Shell: "/bin/bash",
	}
	assert.EqualValues(t, expected, users[0])
}

func TestReadDBFromFileInvalidFilePath(t *testing.T) {
	users, err := ReadDBFromFile(1000, 2000, 1000, 2000, "/foobar/yadayada")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such file")
	assert.Equalf(t, 0, len(users), "Unexpected content: %v", users)
}
