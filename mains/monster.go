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
    var fd *os.File

    arguments(&opts)
    opts.prodfile = flag.Args()[0]
    if opts.prodfile == ""  ||  opts.help {
        usage()
        os.Exit(1)
    }

    popts.Prodfile = opts.prodfile
    popts.Rnd = rand.New( rand.NewSource( int64(opts.seed) ))
    start := monster.Parse( &popts )
    fmt.Printf("start - %v\n", start)

    if opts.ast {
        monster.PrintAst( popts.Nonterminals )
    } else {
        var s string
        nt := popts.Nonterminals[ start ]
        if opts.outfile != "-" {
            fd, err := os.Create( opts.outfile )
            if err != nil {
                fmt.Printf("%v\n", err )
                os.Exit(1)
            }
            defer func() { if err = fd.Close(); err != nil { panic(err) } }()
        } else {
            fd = os.Stdout
        }
        for i:=0; i < opts.count; i++ {
            s = monster.Generate(monster.Context{}, popts, nt)
            fd.Write( bytes.NewBufferString(s).Bytes() )
            fd.Write( bytes.NewBufferString("\n\n").Bytes() )
        }
    }
}
