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
		UID:  1000,
		GID:  1000,
		GECOS: GECOS{
			FullName:  "Alceu Rodrigues de Freitas Junior",
			Office:    "571",
			WorkPhone: "+551155422748",
			HomePhone: "+551155422748",
			Raw:       "Alceu Rodrigues de Freitas Junior,571,+551155422748,+551155422748",
		},
		HomeDir: "/home/alceu",
		Shell:   "/bin/bash",
	}
	assert.EqualValues(t, expected, users[0])
}

func TestReadDBFromFileInvalidFilePath(t *testing.T) {
	users, err := ReadDBFromFile(1000, 2000, 1000, 2000, "/foobar/yadayada")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no such file")
	assert.Equalf(t, 0, len(users), "Unexpected content: %v", users)
}

func TestToAccountLDIF(t *testing.T) {
	users, err := ReadDBFromFile(1000, 2000, 1000, 2000, mockDBFilename)
	assert.Nilf(t, err, "Unexpected error: '%s'", err)
	expected := []string{
		"dn: uid=alceu,ou=People,dc=foobar,dc=org",
		"uid: alceu",
		"loginShell: /bin/bash",
		"uidNumber: 1000",
		"gidNumber: 1000",
		"homeDirectory: /home/alceu",
		"cn: Alceu Rodrigues de Freitas Junior",
		"objectClass: posixAccount",
		"objectClass: top",
		"objectClass: account",
		"gecos: Alceu Rodrigues de Freitas Junior,571,+551155422748,+551155422748",
	}
	assert.EqualValues(t, expected, users[0].ToAccountLDIF("dc=foobar,dc=org"))
}

func TestToPersonLDIF(t *testing.T) {
	users, err := ReadDBFromFile(1000, 2000, 1000, 2000, mockDBFilename)
	assert.Nilf(t, err, "Unexpected error: '%s'", err)
	expected := []string{
		"dn: uid=alceu,ou=People,dc=foobar,dc=org",
		"uid: alceu",
		"mail: alceu@foobar.org",
		"telephoneNumber: +551155422748",
		"roomNumber: 571",
		"homePhone: +551155422748",
		"givenName: Alceu",
		"sn: Rodrigues de Freitas Junior",
		"cn: Alceu Rodrigues de Freitas Junior",
		"objectClass: posixAccount",
		"objectClass: top",
		"objectClass: person",
		"objectClass: organizationalPerson",
		"objectClass: inetOrgPerson",
		"loginShell: /bin/bash",
		"uidNumber: 1000",
		"gidNumber: 1000",
		"homeDirectory: /home/alceu",
	}
	assert.EqualValues(t, expected, users[0].ToPersonLDIF("foobar.org", "", "dc=foobar,dc=org"))
}

func TestToPersonLDIFWithMailHost(t *testing.T) {
	users, err := ReadDBFromFile(1000, 2000, 1000, 2000, mockDBFilename)
	assert.Nilf(t, err, "Unexpected error: '%s'", err)
	expected := []string{
		"dn: uid=alceu,ou=People,dc=foobar,dc=org",
		"uid: alceu",
		"mail: alceu@foobar.org",
		"telephoneNumber: +551155422748",
		"roomNumber: 571",
		"homePhone: +551155422748",
		"givenName: Alceu",
		"sn: Rodrigues de Freitas Junior",
		"cn: Alceu Rodrigues de Freitas Junior",
		"mailRoutingAddress: alceu@mailer.foobar.org",
		"mailHost: mailer.foobar.org",
		"objectClass: inetLocalMailRecipient",
		"objectClass: posixAccount",
		"objectClass: top",
		"objectClass: person",
		"objectClass: organizationalPerson",
		"objectClass: inetOrgPerson",
		"loginShell: /bin/bash",
		"uidNumber: 1000",
		"gidNumber: 1000",
		"homeDirectory: /home/alceu",
	}
	assert.EqualValues(t, expected, users[0].ToPersonLDIF("foobar.org", "mailer.foobar.org", "dc=foobar,dc=org"))
}
