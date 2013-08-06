package monster
import ("fmt"; "os"; "strconv"; "math/rand"; "encoding/csv")

var Bagfiles = make( map[string][][]string )    // A cache of bag files.

func bag( c Context, nt *NonTerminal ) string {
    var filename string
    var index int

    cs := nt.Children // Arguments
    if len(cs) == 2 {
        filename = cs[0].(*Terminal).Value
        index, _ = strconv.Atoi(cs[1].(*Terminal).Value)
    } else if len(cs) == 1 {
        filename, index = cs[0].(*StrTerminal).Value, 0
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
    BnfCallbacks["bag"] = bag
}
