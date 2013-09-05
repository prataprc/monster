package monster

import (
    "fmt"
    "os"
    "strconv"
    "math/rand"
    "encoding/csv"
)

var Bagfiles = make(map[string][][]string) // A cache of bag files.

func bag(c Context, args []interface{}) string {
    var filename string
    var index int

    if len(args) == 2 {
        filename = args[0].(string)
        index, _ = strconv.Atoi(args[1].(string))
    } else if len(args) == 1 {
        filename, index = args[0].(string), 0
    } else {
        panic("Error: Atleast one argument expected in bag() BNF")
    }
    filename = filename[1: len(filename)-1] // remove the double quotes
    rnd := c["_random"].(*rand.Rand)
    return rangeOnFile(rnd, filename, index)
}

func rangeOnFile(rnd *rand.Rand, filename string, index int) string {
    var choice = Bagfiles[filename]
    if choice == nil {
        choice = readBag( filename )
        Bagfiles[filename] = choice
    }
    record := choice[ rnd.Intn(len(choice)) ]
    return record[index]
}

func readBag(filename string) [][]string {
    fd, err := os.Open(filename)
    if err != nil {
        fmt.Printf( "Cannot open file %v\n", filename )
        os.Exit(1)
    }
    records, _ := csv.NewReader(fd).ReadAll()
    return records
}

func init() {
    BnfCallbacks["bag"] = bag
}
