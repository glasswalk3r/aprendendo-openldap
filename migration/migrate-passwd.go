package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"migration.openldap.org/passwd/db/shadow"
	"migration.openldap.org/passwd/db/passwd"
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

	flag.StringVar(&dnsDomain, "dns-domain", defaultDNSDomain, fmt.Sprintf("Specify the DNS domain to use, default to %s", defaultDNSDomain))
	flag.StringVar(&baseDN, "base-dn", defaultBaseDN, fmt.Sprintf("Specify the base DN, default to %s", defaultBaseDN))
	flag.StringVar(&mailHost, "mail-host", "", "Optional, define inetLocalMailRecipient information if available")
	flag.IntVar(&uidBelow, "ignore-uid-below", defaultBelow, fmt.Sprintf("Specify the minimum UID to consider retrieving, default is %d", defaultBelow))
	flag.IntVar(&uidAbove, "ignore-uid-above", defaultAbove, fmt.Sprintf("Specify the maximum UID to consider retrieving, default is %d", defaultAbove))
	flag.IntVar(&gidBelow, "ignore-gid-below", defaultBelow, fmt.Sprintf("Specify the minimum GID to consider retrieving, default is %d", defaultBelow))
	flag.IntVar(&gidAbove, "ignore-gid-above", defaultAbove, fmt.Sprintf("Specify the maximum GID to consider retrieving, default is %d", defaultAbove))
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

	for _, entry := range passwdDB {
		dump := entry.ToLDIF(dnsDomain, mailHost)
		fmt.Println(strings.Join(dump, "\n"))
		shadowEntry, err := shadowDB.UserEntry(entry.User)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		fmt.Println(strings.Join(shadowEntry.ToLDIF(), "\n"))
		fmt.Println()
	}
}
