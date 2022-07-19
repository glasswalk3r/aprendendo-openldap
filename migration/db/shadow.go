package shadow

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

func NewDB(filePath string) ShadowDB {

	// basically to enable unit testing
	if filePath == "" {
		filePath = ShadowFilePath
	}

	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read %s: %s", filePath, err)
		os.Exit(1)
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

	return db
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
