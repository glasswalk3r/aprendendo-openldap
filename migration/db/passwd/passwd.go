// Provides features to recover information from /etc/passwd file in Linux
package passwd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
DBEntry is the representation of a user in the /etc/passwd file.
The password field is ignored since modern UNIX-like OS uses /etc/shadow for
storing a hashing of the original password
*/
type DBEntry struct {
	User    string
	UID     int
	GID     int
	GECOS   GECOS
	HomeDir string
	Shell   string
}

/*
GECOS is an arbitrary list of string separated by commas.
See https://www.redhat.com/sysadmin/linux-gecos-demystified for more details.
*/
type GECOS struct {
	FullName  string
	Office    string
	WorkPhone string
	HomePhone string
	Raw       string
}

type PersonName struct {
	GivenName string
	Surname   string
}

// Splits a GECOS FullName between given name and surname.
func (g *GECOS) SplitName() PersonName {
	if g.FullName == "" {
		return PersonName{}
	}

	parts := strings.Split(g.FullName, " ")
	return PersonName{
		parts[0],
		strings.Join(parts[1:], " "),
	}
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

	return DBEntry{
		User:    user[0],
		UID:     int(uid),
		GID:     int(gid),
		GECOS:   gecos,
		HomeDir: user[5],
		Shell:   user[6],
	}, nil
}

/*
ToAccountLDIF exports a DBEntry struct in a LDIF format with attributes for the
account and posixAccount classes.
*/
func (e *DBEntry) ToAccountLDIF(baseDN string) []string {
	var cn string

	if e.GECOS.FullName != "" {
		cn = fmt.Sprintf("cn: %s", e.GECOS.FullName)
	} else {
		cn = fmt.Sprintf("cn: %s", e.User)
	}

	dump := []string{
		fmt.Sprintf("dn: uid=%s,ou=People,%s", e.User, baseDN),
		fmt.Sprintf("uid: %s", e.User),
		fmt.Sprintf("loginShell: %s", e.Shell),
		fmt.Sprintf("uidNumber: %d", e.UID),
		fmt.Sprintf("gidNumber: %d", e.GID),
		fmt.Sprintf("homeDirectory: %s", e.HomeDir),
		cn,
	}

	objectClasses := []string{"posixAccount", "top", "account"}

	for _, value := range objectClasses {
		dump = append(dump, fmt.Sprintf("objectClass: %s", value))
	}

	if e.GECOS.Raw != "" {
		dump = append(dump, fmt.Sprintf("gecos: %s", e.GECOS.Raw))
	}

	return dump
}

/*
ToPersonLDIF exports a DBEntry struct in a LDIF format with attributes for
person, organizationalPerson and inetOrgPerson classes.
Expect as parameters:
- dnsDomain: specify the DNS domain to use with the mail attribute
- mailHost: define inetLocalMailRecipient class attributes
- baseDN: specify the base DN for the entry DN
*/
func (e *DBEntry) ToPersonLDIF(dnsDomain, mailHost, baseDN string) []string {
	// TODO: try refactoring by using append() instead
	dump := []string{
		fmt.Sprintf("dn: uid=%s,ou=People,%s", e.User, baseDN),
		fmt.Sprintf("uid: %s", e.User),
		fmt.Sprintf("mail: %s@%s", e.User, dnsDomain),
	}

	if e.GECOS.Raw != "" {
		if e.GECOS.WorkPhone != "" {
			dump = append(dump, fmt.Sprintf("telephoneNumber: %s", e.GECOS.WorkPhone))
		}

		if e.GECOS.Office != "" {
			dump = append(dump, fmt.Sprintf("roomNumber: %s", e.GECOS.Office))
		}

		if e.GECOS.HomePhone != "" {
			dump = append(dump, fmt.Sprintf("homePhone: %s", e.GECOS.HomePhone))
		}

		if e.GECOS.FullName != "" {
			pn := e.GECOS.SplitName()
			dump = append(dump, fmt.Sprintf("givenName: %s", pn.GivenName))
			dump = append(dump, fmt.Sprintf("sn: %s", pn.Surname))
			dump = append(dump, fmt.Sprintf("cn: %s", e.GECOS.FullName))
		}
	} else {
		dump = append(dump, fmt.Sprintf("cn: %s", e.User))
	}

	if mailHost != "" {
		dump = append(dump, fmt.Sprintf("mailRoutingAddress: %s@%s", e.User, mailHost))
		dump = append(dump, fmt.Sprintf("mailHost: %s", mailHost))
		dump = append(dump, "objectClass: inetLocalMailRecipient")
	}

	objectClasses := []string{"posixAccount", "top", "person", "organizationalPerson", "inetOrgPerson"}

	for _, value := range objectClasses {
		dump = append(dump, fmt.Sprintf("objectClass: %s", value))
	}

	dump = append(dump, fmt.Sprintf("loginShell: %s", e.Shell))
	dump = append(dump, fmt.Sprintf("uidNumber: %d", e.UID))
	dump = append(dump, fmt.Sprintf("gidNumber: %d", e.GID))
	dump = append(dump, fmt.Sprintf("homeDirectory: %s", e.HomeDir))
	return dump
}

/*
NewGECOS is the constructor for the GECOS struct.
Only the first four fields considered when parsing, but the original field is
also kept (Raw).
Expects the GECOS field as available in /etc/passwd.
*/
func NewGECOS(gecos string) GECOS {
	current := strings.Split(gecos, ",")
	expected := [4]string{}

	for i, _ := range current {
		expected[i] = current[i]
	}

	return GECOS{
		FullName:  expected[0],
		Office:    expected[1],
		WorkPhone: expected[2],
		HomePhone: expected[3],
		Raw:       gecos,
	}
}

/*
ReadDB reads all the users from the /etc/passwd and return those which
respective UID and GID pass the provided filters.
The mingGID parameter is the minimum GID number that will be considered to
retrieve, meaning that GID's lesser than that will be ignored.
The maxGID is used for the same GID filter, but GID's greater than the specified
value will be ignored.
Unless you're doing unit testing, this is the function you should be using.
*/
func ReadDB(minUID, maxUID, minGID, maxGID int) ([]DBEntry, error) {
	return ReadDBFromFile(minUID, maxUID, minGID, maxGID, "/etc/passwd")
}

/*
ReadDBFromFile does the same thing as ReadDB, but reads from an arbitrary
file location, which is good for unit testing.
*/
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
