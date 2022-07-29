package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"migration.openldap.org/passwd/db/group"
	"migration.openldap.org/passwd/db/passwd"
	"migration.openldap.org/passwd/db/shadow"
)

const defaultDNSDomain string = "foobar.org"
const defaultBaseDN string = "dc=foobar,dc=org"
const defaultBelow int = 1000
const defaultAbove int = 2000

func main() {
	var dnsDomain string
	var baseDN string
	var mailHost string
	var uidBelow int
	var uidAbove int
	var gidBelow int
	var gidAbove int
	var writeResultTo string
	var useExtended bool

	flag.StringVar(&dnsDomain, "dns-domain", defaultDNSDomain, fmt.Sprintf("Specify the DNS domain to use, default to %s", defaultDNSDomain))
	flag.StringVar(&baseDN, "base-dn", defaultBaseDN, fmt.Sprintf("Specify the base DN, default to %s", defaultBaseDN))
	flag.StringVar(&mailHost, "mail-host", "", "Optional, define inetLocalMailRecipient attributes if provided")
	flag.StringVar(&writeResultTo, "save-to", "", "Optional, path to a file to save LDIF result if provided")
	flag.IntVar(&uidBelow, "ignore-uid-below", defaultBelow, fmt.Sprintf("Specify the minimum UID to consider retrieving, default is %d", defaultBelow))
	flag.IntVar(&uidAbove, "ignore-uid-above", defaultAbove, fmt.Sprintf("Specify the maximum UID to consider retrieving, default is %d", defaultAbove))
	flag.IntVar(&gidBelow, "ignore-gid-below", defaultBelow, fmt.Sprintf("Specify the minimum GID to consider retrieving, default is %d", defaultBelow))
	flag.IntVar(&gidAbove, "ignore-gid-above", defaultAbove, fmt.Sprintf("Specify the maximum GID to consider retrieving, default is %d", defaultAbove))
	flag.BoolVar(&useExtended, "use-extended", false, "Uses the LDAP inetOrgPerson class for extended attributes, otherwise Account will be used by default")
	flag.Parse()

	shadowDB, err := shadow.ReadDB()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	passwdDB, err := passwd.ReadDB(uidBelow, uidAbove, gidBelow, gidAbove)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var writeFile *os.File
	var fileWriter *bufio.Writer

	if writeResultTo != "" {
		writeFile, err = os.Create(writeResultTo)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		defer writeFile.Close()
		fileWriter = bufio.NewWriter(writeFile)
	}

	for _, entry := range passwdDB {
		var dump []string

		if useExtended {
			dump = entry.ToPersonLDIF(dnsDomain, mailHost, baseDN)
		} else {
			dump = entry.ToAccountLDIF(baseDN)
		}

		if writeResultTo != "" {
			fileWriter.WriteString(strings.Join(dump, "\n"))
			fileWriter.WriteString("\n")
		} else {
			fmt.Println(strings.Join(dump, "\n"))
		}

		shadowEntry, err := shadowDB.UserEntry(entry.User)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		if writeResultTo != "" {
			fileWriter.WriteString(strings.Join(shadowEntry.ToLDIF(), "\n"))
			// put a required new line between two different entries
			fileWriter.WriteString("\n\n")
			fileWriter.Flush()
		} else {
			fmt.Println(strings.Join(shadowEntry.ToLDIF(), "\n"))
			fmt.Println()
		}
	}

	groups, err := group.ReadDB(gidBelow, gidAbove)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if writeResultTo != "" {
		for _, g := range groups {
			fileWriter.WriteString(strings.Join(g.ToLDIF(dnsDomain, mailHost, baseDN), "\n"))
			fileWriter.WriteString("\n\n")
		}

		fileWriter.Flush()
	} else {
		for _, g := range groups {
			fmt.Println(strings.Join(g.ToLDIF(dnsDomain, mailHost, baseDN), "\n"))
			fmt.Println()
		}
	}
}
