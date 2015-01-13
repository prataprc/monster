//  Copyright (c) 2013 Couchbase, Inc.

package builtin

import "fmt"
import "os"
import "encoding/csv"
import "path/filepath"

import "github.com/prataprc/monster/common"

var cacheBagRecords = make(map[string][][]string)

// Bag will fetch a random line from file and return it.
// args[0] - filename.
func Bag(scope common.Scope, args ...interface{}) interface{} {
	var err error

	filename := args[0].(string)
	if !filepath.IsAbs(filename) {
		if bagdir, _, ok := scope.GetString("_bagdir"); ok {
			filename = filepath.Join(bagdir, filename)
		} else if prodfile, _, ok := scope.GetString("_prodfile"); ok {
			dirpath := filepath.Dir(prodfile)
			filename = filepath.Join(dirpath, filename)
		}
	}
	if filename, err = filepath.Abs(filename); err != nil {
		panic(fmt.Errorf("bad filepath: %v\n", filename))
	}

	records, ok := cacheBagRecords[filename]
	if !ok {
		records = readBag(filename)
		cacheBagRecords[filename] = records
	}
	if len(records) > 0 {
		rnd := scope.GetRandom()
		record := records[rnd.Intn(len(records))]
		if len(record) > 0 {
			return record[0]
		}
	}
	return ""
}

func readBag(filename string) [][]string {
	fd, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("cannot open file %v\n", filename))
	}
	records, err := csv.NewReader(fd).ReadAll()
	if err == nil {
		return records
	}
	panic(fmt.Errorf("unable to read file %v in CSV format\n", filename))
}
