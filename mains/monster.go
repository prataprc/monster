package main
import "fmt"
import "flag"
import "math/rand"
import "time"
import "os"
import "bytes"
import "github.com/prataprc/monster"

type options struct {
    ast bool
    prodfile string
    outfile string
    seed int
    count int
    help bool
}

func arguments( opts *options ) {
    seed := time.Now().UTC().Second()
    flag.BoolVar( &opts.ast, "ast", false, "Show the ast of production" )
    flag.IntVar( &opts.seed, "s", seed, "Seed value" )
    flag.IntVar( &opts.count, "n", 1, "Generate n combinations" )
    flag.StringVar( &opts.outfile, "o", "-", "Specify an output file" )
    flag.BoolVar( &opts.help, "h", false, "Print usage and default options" )
    flag.Parse()
}

func usage() {
    fmt.Fprintf(os.Stderr, "Usage : %s [OPTIONS] <production-file> \n", os.Args[0])
    flag.PrintDefaults()
}

func main() {
    var opts = options{}
    var popts = monster.ParseOpts{}

    arguments(&opts)
    opts.prodfile = flag.Args()[0]
    if opts.prodfile == ""  ||  opts.help {
        usage()
        os.Exit(1)
    }

    popts.Prodfile = opts.prodfile
    popts.Rnd = rand.New( rand.NewSource( int64(opts.seed) ))
    start, nonterminals := monster.Parse( popts )
    fmt.Printf("start - %v\n", start)

    if opts.ast {
        monster.PrintAst( nonterminals )
    } else {
        nt := nonterminals[ start ]
        if opts.outfile != "-" {
            fd, err := os.Create( opts.outfile )
            if err != nil {
                fmt.Printf("%v\n", err )
                os.Exit(1)
            }
            defer func() { if err = fd.Close(); err != nil { panic(err) } }()
            out( nt, opts.count, popts, fd )
        } else {
            out( nt, opts.count, popts, os.Stdout )
        }
    }
}

func out( nt monster.NonTerminal, count int, popts monster.ParseOpts, fd *os.File ) {
    var s string
    for i:=0; i < count; i++ {
        s = monster.Generate(popts, nt)
        fd.Write( bytes.NewBufferString(s).Bytes() )
        fd.Write( bytes.NewBufferString("\n\n").Bytes() )
    }
}
