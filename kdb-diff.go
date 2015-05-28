package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type Database struct {
	Groups []Group `xml:"group"`
}

type Group struct {
	Title   string  `xml:"title"`
	Entries []Entry `xml:"entry"`
	Groups  []Group `xml:"group"`
}

type Entry struct {
	Title    string `xml:"title"`
	Username string `xml:"username"`
	Password string `xml:"password"`
	Url      string `xml:"url"`
	Comment  string `xml:"comment"`
}

func filterEntries(groups []Group, entries *map[string]Entry) {
	for _, group := range groups {
		if group.Title == "Backup" {
			fmt.Println("Except Backup.")
			return
		}
		if len(group.Groups) > 0 {
			filterEntries(group.Groups, entries)
		}
		for _, entry := range group.Entries {
			//fmt.Println(entry.Title)
			if duplicatedEntry, ok := (*entries)[entry.Title]; ok {
				fmt.Println("duplicated: ", duplicatedEntry.Title)
			} else {
				(*entries)[entry.Title] = entry
			}
		}
	}
}

func diff(masterEntries map[string]Entry, conflictEntries map[string]Entry) {
	masterOnlyKeys := make([]string, 0)
	conflictOnlyKeys := make([]string, 0)
	diffKeys := make([]string, 0)

	for k, v := range masterEntries {
		if conflictEntry, ok := conflictEntries[k]; ok {
			isDiff := false
			if v.Title != conflictEntry.Title {
				isDiff = true
			} else if v.Username != conflictEntry.Username {
				isDiff = true
			} else if v.Password != conflictEntry.Password {
				isDiff = true
			} else if v.Url != conflictEntry.Url {
				isDiff = true
			} else if v.Comment != conflictEntry.Comment {
				isDiff = true
			}
			if isDiff {
				diffKeys = append(diffKeys, k)
			}
		} else {
			masterOnlyKeys = append(masterOnlyKeys, k)
		}
	}
	for k, _ := range conflictEntries {
		if _, ok := masterEntries[k]; !ok {
			conflictOnlyKeys = append(conflictOnlyKeys, k)
		}
	}

	fmt.Println("Master Only: ", masterOnlyKeys)
	fmt.Println("Conflict Only: ", conflictOnlyKeys)
	fmt.Println("Diff: ", diffKeys)
}

func readXML(filePath string) (d Database, err error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return d, err
	}
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(XMLdata, &d)
	return d, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: %s FILE1 FILE2\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	filePaths := flag.Args()

	masterDB, err := readXML(filePaths[0])
	if err != nil {
		return
	}
	conflictDB, err := readXML(filePaths[1])
	if err != nil {
		return
	}
	masterEntries := make(map[string]Entry)
	conflictEntries := make(map[string]Entry)
	filterEntries(masterDB.Groups, &masterEntries)
	filterEntries(conflictDB.Groups, &conflictEntries)
	diff(masterEntries, conflictEntries)
}
