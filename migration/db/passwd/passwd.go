package passwd

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// DBEntry is the representation of a user in the /etc/passwd file.
// The password field is ignored since modern UNIX-like OS uses /etc/shadow for
// storing a hashing of the original password
type DBEntry struct {
	User    string
	UID     int
	GID     int
	GECOS   GECOS
	Homedir string
	Shell   string
}

// GECOS is an arbitrary list of string separated by commas.
// See https://www.redhat.com/sysadmin/linux-gecos-demystified for more details.
type GECOS struct {
	FullName  string
	Office    string
	WorkPhone string
	HomePhone string
}

// NewDBEntry is the constructor for the DBEntry struct
func NewDBEntry(user []string) (DBEntry, error) {
	gecos := NewGECOS(user[4])
	uid, err := strconv.ParseInt(user[2], 0, 0)

	if err != nil {
		return DBEntry{}, err
	}

	gid, err := strconv.ParseInt(user[3], 0, 0)

	if err != nil {
		return DBEntry{}, err
	}

	return DBEntry{user[0], int(uid), int(gid), gecos, user[5], user[6]}, nil
}

// NewGECOS is the constructor for the GECOS struct
func NewGECOS(gecos string) GECOS {
	current := strings.Split(gecos, ",")
	expected := [4]string{}

	for i, _ := range current {
		expected[i] = current[i]
	}

	return GECOS{expected[0], expected[1], expected[2], expected[3]}
}

// ReadDB reads all the users from the /etc/passwd and return those
// UID and GID pass the provided filters.
// Unless you're doing unit testing, this is the function you should be using
// to start of
func ReadDB(minUID, maxUID, minGID, maxGID int) ([]DBEntry, error) {
  return ReadDBFromFile(minUID, maxUID, minGID, maxGID, "/etc/passwd")
}

// ReadDBFromFile does the same thing as ReadDB, but reads from an arbitrary
// file location, which is good for unit testing
func ReadDBFromFile(minUID, maxUID, minGID, maxGID int, filePath string) ([]DBEntry, error) {
	var users []DBEntry
	readFile, err := os.Open(filePath)

	if err != nil {
		return users, err
	}

	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		entry := strings.Split(line, ":")
		user, err := NewDBEntry(entry)

		if err != nil {
			return users, err
		}

		if user.UID >= minUID && user.UID <= maxUID && user.GID >= minGID && user.GID <= maxGID {
			users = append(users, user)
		}
	}

	return users, nil
}
