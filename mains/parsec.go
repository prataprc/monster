package main
import "fmt"
import "os"
import "flag"
import "github.com/prataprc/monster/parsec"

type options struct {
    grammarfile string
}

func arguments( opts *options ) {
    flag.Parse()
}

func usage() {
    fmt.Fprintf(os.Stderr, "Usage : %s [OPTIONS] <grammar-file> \n", os.Args[0])
    flag.PrintDefaults()
}

func main() {
    var opts = options{}
    //var fd *os.File

    arguments(&opts)
    opts.grammarfile = flag.Args()[0]
    if opts.grammarfile == ""  {
        usage()
        os.Exit(1)
    }
    scanner := parsec.NewGoScanner( opts.grammarfile )
    for {
        tok := scanner.Scan()
        fmt.Printf("%v \n", tok)
        if tok.Type == "EOF" { break }
    }
}
