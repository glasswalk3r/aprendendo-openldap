package db

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"errors"
)

type shadowDBAttribs [8]string

type ShadowDB struct {
	rawEntries map[string]shadowDBAttribs
}

type ShadowDBEntry struct {
	userPassword     string
	shadowLastChange string
	shadowMin        string
	shadowMax        string
	shadowWarning    string
	shadowInactive   string
	shadowExpire     string
	shadowFlag       string
}

const ShadowFilePath string = "/etc/shadow"
const ShadowFieldSep string = ":"

// ReadDB is basically a shortcut to NewDBFromFile, using by default the standard path to the shadow file (see ShadowFilePath constant)
func ReadDB() (ShadowDB, error) {
	return ReadDBFromFile(ShadowFilePath)
}

// ReadDBFromFile reads a text file that has lines in the format documented in shadow section 5 manpage
func ReadDBFromFile(filePath string) (ShadowDB, error) {
	readFile, err := os.Open(filePath)

	if err != nil {
		return ShadowDB{}, err
	}

	defer readFile.Close()
	var db ShadowDB
	db.rawEntries = make(map[string]shadowDBAttribs)

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		entry := fileScanner.Text()
		fields := strings.Split(entry, ShadowFieldSep)
		db.rawEntries[fields[0]] = shadowDBAttribs{fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], fields[7], fields[8]}
	}

	return db, nil
}

func (db *ShadowDB) UserEntry(user string) (ShadowDBEntry, error) {

	if user == "" {
		return ShadowDBEntry{}, errors.New("user is a required parameter")
	}

	attribs, ok := db.rawEntries[user]

	if !ok {
		return ShadowDBEntry{}, fmt.Errorf("user '%s' not found in the database", user)
	}

	return ShadowDBEntry{
		attribs[0],
		attribs[1],
		attribs[2],
		attribs[3],
		attribs[4],
		attribs[5],
		attribs[6],
		attribs[7],
	}, nil
}

func (e *ShadowDBEntry) ToLDIF() []string {
	// 1 + attributes
	var dump [9]string
	dump[0] = "objectClass: shadowAccount"
	lastInUse := 1

	if e.userPassword != "" {
		dump[1] = fmt.Sprintf("userPassword: {crypt}%s", e.userPassword)
		lastInUse++
	}

	if e.shadowLastChange != "" {
		dump[2] = fmt.Sprintf("shadowLastChange: %s", e.shadowLastChange)
		lastInUse++
	}

	if e.shadowMin != "" {
		dump[3] = fmt.Sprintf("shadowMin: %s", e.shadowMin)
		lastInUse++
	}

	if e.shadowMax != "" {
		dump[4] = fmt.Sprintf("shadowMax: %s", e.shadowMax)
		lastInUse++
	}

	if e.shadowWarning != "" {
		dump[5] = fmt.Sprintf("shadowWarning: %s", e.shadowWarning)
		lastInUse++
	}

	if e.shadowInactive != "" {
		dump[6] = fmt.Sprintf("shadowInactive: %s", e.shadowInactive)
		lastInUse++
	}

	if e.shadowExpire != "" {
		dump[7] = fmt.Sprintf("shadowExpire: %s", e.shadowExpire)
		lastInUse++
	}

	if e.shadowFlag != "" {
		dump[8] = fmt.Sprintf("shadowFlag: %s", e.shadowFlag)
		lastInUse++
	}

	return dump[:lastInUse]
}
