//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package monster

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

// Bagfiles stores a cache of bag entries to avoid loading bag files.
var Bagfiles = make(map[string][][]string)

func bag(c Context, args []interface{}) string {
	var filename string
	var index int
	var err error

	if len(args) == 2 {
		filename = args[0].(string)
		index, _ = strconv.Atoi(args[1].(string))
	} else if len(args) == 1 {
		filename, index = args[0].(string), 0
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
	rnd := c["_random"].(*rand.Rand)
	return rangeOnFile(rnd, filename, index)
}

func rangeOnFile(rnd *rand.Rand, filename string, index int) string {
	var choice = Bagfiles[filename]
	if choice == nil {
		choice = readBag(filename)
		Bagfiles[filename] = choice
	}
	record := choice[rnd.Intn(len(choice))]
	return record[index]
}

func readBag(filename string) [][]string {
	fd, err := os.Open(filename)
	if err != nil {
		panic(fmt.Errorf("cannot open file %v\n", filename))
	}
	records, _ := csv.NewReader(fd).ReadAll()
	return records
}

func init() {
	BnfCallbacks["bag"] = bag
}
