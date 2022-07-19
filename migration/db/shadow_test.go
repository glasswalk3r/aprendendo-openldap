package shadow

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const mockDBFilename string = "mock.txt"
const userEntry string = "johndoe"
const userPassword string = "uselesspassword"

func TestNewDB(t *testing.T) {
	db := NewDB(mockDBFilename)
	assert.Equal(t, 1, len(db.rawEntries))
	attribs, ok := db.rawEntries[userEntry]
	assert.Truef(t, ok, "User entry '%s' not found", userEntry)
	assert.Equal(t, userPassword, attribs[0])
}

func TestUserEntry(t *testing.T) {
	db := NewDB(mockDBFilename)
	expected := ShadowDBEntry{"uselesspassword", "18729", "0", "99999", "7", "", "", ""}
	entry, err := db.UserEntry(userEntry)
	assert.Nilf(t, err, "Error when trying to retrieve '%s' entry", userEntry)
	assert.EqualValuesf(t, expected, entry, "%v not equal to %v", entry, expected)
}
