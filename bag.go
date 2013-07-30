package monster

import "fmt"
import "os"
import "strconv"
import "encoding/csv"

var Bagfiles = make( map[string][][]string )    // A cache of bag files.

func bag( popts ParseOpts ) Terminal {
    var filename string
    var index int = 0
    var t Terminal
    _ = Token(popts.S) // "("
    toktype, filename := Tokent(popts.S)
    if toktype != "String" {
        fmt.Printf("Error: bag() first argument should be string\n")
        os.Exit(1)
    }
    filename = filename[1: len(filename)-1] // remove the double quotes
    tok := Token(popts.S)
    if tok == "," {
        toktype, tok = Tokent(popts.S)
        if toktype != "Int" {
            fmt.Printf("Error: bag() second argument should be integer\n")
            os.Exit(1)
        }
        index, _ = strconv.Atoi( tok )
        tok = Token(popts.S)
    }
    if tok == ")" {
        fn := func(contex Context) string {
            return rangeOnFile(popts, filename, index)
        }
        return Terminal{ name : "BnfBag", value : "", generator: fn }
    } else {
        fmt.Printf("Syntax error in bnf_bag\n")
        os.Exit(1)
    }
    return t // Dummy return, otherwise compiler cribs
}

func rangeOnFile(popts ParseOpts, filename string, index int) string {
    var choice = Bagfiles[filename]
    if choice == nil {
        choice = readBag( filename )
        Bagfiles[filename] = choice
    }
    record := choice[ popts.Rnd.Intn(len(choice)) ]
    return record[index]
}

func readBag( filename string ) [][]string {
    fd, err := os.Open(filename)
    if err != nil {
        fmt.Printf( "Cannot open file %v\n", filename )
        os.Exit(1)
    }
    records, _ := csv.NewReader(fd).ReadAll()
    return records
}

func init() {
    Bnfs["bag"] = bag
}
