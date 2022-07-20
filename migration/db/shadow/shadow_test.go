package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const mockDBFilename string = "mock.txt"
const userEntry string = "johndoe"
const userPassword string = "uselesspassword"

func TestReadDBFromFile(t *testing.T) {
	db, err := ReadDBFromFile(mockDBFilename)
	assert.Nilf(t, err, "Unexpected error: '%s'", err)
	assert.Equal(t, 1, len(db.rawEntries))
	attribs, ok := db.rawEntries[userEntry]
	assert.Truef(t, ok, "User entry '%s' not found", userEntry)
	assert.Equal(t, userPassword, attribs[0])
}

func TestReadDBFromFileInvalidFilePath(t *testing.T) {
	db, err := ReadDBFromFile("/foobar/foo/bar")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such file")
	assert.Equalf(t, 0, len(db.rawEntries), "Unexpected content: %v", db)
}

func TestUserEntry(t *testing.T) {
	db, err := ReadDBFromFile(mockDBFilename)
	assert.Nilf(t, err, "Unexpected error: '%s'", err)
	expected := ShadowDBEntry{"uselesspassword", "18729", "0", "99999", "7", "", "", ""}
	entry, err := db.UserEntry(userEntry)
	assert.Nilf(t, err, "Error when trying to retrieve '%s' entry", userEntry)
	assert.EqualValuesf(t, expected, entry, "%v not equal to %v", entry, expected)
}

func TestUserEntryNotFound(t *testing.T) {
	db, err := ReadDBFromFile(mockDBFilename)
	entry, err := db.UserEntry("yadayada")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not found")
	expected := ShadowDBEntry{}
	assert.EqualValuesf(t, expected, entry, "%v not equal to %v", entry, expected)
}

func TestToLDIF(t *testing.T) {
	db, _ := ReadDBFromFile(mockDBFilename)
	entry, _ := db.UserEntry(userEntry)
	expected := []string{
		"objectClass: shadowAccount",
		"userPassword: {crypt}uselesspassword",
		"shadowLastChange: 18729",
		"shadowMin: 0",
		"shadowMax: 99999",
		"shadowWarning: 7",
	}
	assert.EqualValues(t, expected, entry.ToLDIF())
}
