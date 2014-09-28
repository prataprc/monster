//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain a
//  copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package monster

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

func bag(c Context, args []interface{}) string {
	var filename string
	var column int
	var err error

	if len(args) == 2 {
		filename = args[0].(string)
		column, _ = strconv.Atoi(args[1].(string))
	} else if len(args) == 1 {
		filename, column = args[0].(string), 0
	} else {
		panic("Error: Atleast one argument expected in bag() BNF")
	}
	filename = filename[1 : len(filename)-1] // remove the double quotes
	if !filepath.IsAbs(filename) {
		if bagdir, ok := c["_bagdir"].(string); ok {
			filename = filepath.Join(bagdir, filename)
		} else if prodfile, ok := c["_prodfile"]; ok {
			dirpath := filepath.Dir(prodfile.(string))
			filename = filepath.Join(dirpath, filename)
		}
	}
	if filename, err = filepath.Abs(filename); err != nil {
		panic("Error: bad filepath")
	}
	return rangeOnFile(c, filename, column)
}

func rangeOnFile(c Context, filename string, column int) string {
	rnd := c["_random"].(*rand.Rand)
	_, ok := c["_bagfiles"]
	if ok == false {
		// Bagfiles stores a cache of bag entries to avoid loading bag files.
		c["_bagfiles"] = make(map[string][][]string)
	}
	bagFiles := c["_bagfiles"].(map[string][][]string)

	choice, ok := bagFiles[filename]
	if ok == false {
		choice = readBag(filename)
		bagFiles[filename] = choice
		return choice[rnd.Intn(len(choice))][column]
	}
	row := rnd.Intn(len(choice))
	return choice[row][column]
}

func readBag(filename string) [][]string {
	fd, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("cannot open file %v\n", filename))
	}
	rs, _ := csv.NewReader(fd).ReadAll()
	records := make([][]string, 0, len(rs))
	for _, record := range rs {
		if len(record) > 0 {
			records = append(records, record)
		}
	}
	return records
}

func init() {
	BnfCallbacks["bag"] = bag
}
