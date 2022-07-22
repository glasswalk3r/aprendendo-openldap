// Provides features to recover information from /etc/group file in Linux
package group

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DBEntry represents a single entry (line) in the /etc/group
type DBEntry struct {
	name     string
	password string
	gid      int
	members  []string
}

// NewDBEntry is the constructor for the DBEntry struct
func NewDBEntry(group []string) (DBEntry, error) {
	gid, err := strconv.ParseInt(group[2], 0, 0)

	if err != nil {
		return DBEntry{}, err
	}

	members := []string{}

	if group[3] != "" {
		members = strings.Split(group[3], ",")
	}

	return DBEntry{
		name:    group[0],
		gid:     int(gid), // password can be safely ignored
		members: members,
	}, nil
}

/*
ReadDB reads all the groups from the /etc/group and return those which the
respective GID pass the provided filters. If required, /etc/gshadow will be read
as well.
The mingGID parameter is the minimum GID number that will be considered to
retrieve, meaning that GID's lesser than that will be ignored.
The maxGID is used for the same GID filter, but GID's greater than the specified
value will be ignored.
Unless you're doing unit testing, this is the function you should be using.
*/
func ReadDB(minGID, maxGID int) ([]DBEntry, error) {
	return ReadDBFromFile(minGID, maxGID, "/etc/group", "/etc/gshadow")
}

/*
ReadDBFromFile does the same thing as ReadDB, but reads from arbitrary files
locations, which is good for unit testing.
*/
func ReadDBFromFile(minGID, maxGID int, groupFilePath, shadowFilePath string) ([]DBEntry, error) {
	var groups []DBEntry
	readGroup, err := os.Open(groupFilePath)

	if err != nil {
		return groups, err
	}

	defer readGroup.Close()
	gshadowDB, err := readShadowDB(shadowFilePath)

	if err != nil {
		return groups, err
	}

	fileScanner := bufio.NewScanner(readGroup)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		entry := strings.Split(line, ":")
		group, err := NewDBEntry(entry)

		if err != nil {
			return groups, err
		}

		if group.gid >= minGID && group.gid <= maxGID {
			password, ok := (*gshadowDB)[group.name]

			if ok {
				group.password = password
			}

			groups = append(groups, group)
		}
	}

	return groups, nil
}

func readShadowDB(shadowFilePath string) (*map[string]string, error) {
	readShadow, err := os.Open(shadowFilePath)
	groups := make(map[string]string)

	if err != nil {
		return &groups, err
	}

	defer readShadow.Close()
	fileScanner := bufio.NewScanner(readShadow)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		fields := strings.Split(line, ":")

		// in most cases the password doesn't exist, so entries can be skipped
		if fields[1] != "!" && fields[1] != "*" {
			// only the password is required
			groups[fields[0]] = fields[1]
		}
	}

	return &groups, nil
}

func (e *DBEntry) ToLDIF(dnsDomain, mailHost, baseDN string) []string {
	dump := make([]string, 5, 10)
	dump[0] = fmt.Sprintf("dn: cn=%s,%s", e.name, baseDN)
	dump[1] = "objectClass: posixGroup"
	dump[2] = "objectClass: top"
	dump[3] = fmt.Sprintf("cn: %s", e.name)
	dump[4] = fmt.Sprintf("gidNumber: %d", e.gid)

	for _, member := range e.members {
		dump = append(dump, fmt.Sprintf("memberUid: %s", member))
	}

	if e.password != "" {
		dump = append(dump, fmt.Sprintf("userPassword: {crypt}%s", e.password))
	}

	return dump
}
