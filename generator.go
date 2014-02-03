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
	"fmt"
	"math/rand"
	"os"
)

// Generate and return an array of text, of `count` length, from `prodfile`.
// `seed` will be used for randomization
// `bagdir` will specify the directory to look for bag-files.
func Generate(seed, count int, bagdir, prodfile string) ([]string, error) {
	if prodfile == "" {
		return nil, fmt.Errorf("invalid prodfile")
	}
	if _, err := os.Stat(prodfile); err != nil {
		return nil, fmt.Errorf("invalid prodfile")
	}
	if _, err := os.Stat(bagdir); bagdir != "" && err != nil {
		return nil, fmt.Errorf("invalid bagdir")
	}

	conf := make(map[string]interface{})
	start, err := Parse(prodfile, conf)
	if err != nil {
		return nil, err
	}

	nonterminals, root := Build(start)
	c := map[string]interface{}{
		"_random":       rand.New(rand.NewSource(int64(seed))),
		"_nonterminals": nonterminals,
		"_bagdir":       bagdir,
	}
	outs := make([]string, 0)
	for i := 0; i < count; i++ {
		Initialize(c)
		outs = append(outs, root.Generate(c))
	}
	return outs, nil
}
