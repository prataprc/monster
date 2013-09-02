package main
import ("fmt"; "flag"; "time"; "os"; "math/rand")
import "github.com/prataprc/monster"
import "github.com/prataprc/golib"

type Interface interface{}
var options struct {
    ast bool
    prodfile string
    outfile string
    seed int
    random *rand.Rand
    count int
    help bool
}

func arguments() {
    seed := time.Now().UTC().Second()
    flag.BoolVar( &options.ast, "ast", false, "Show the ast of production" )
    flag.IntVar( &options.seed, "s", seed, "Seed value" )
    flag.IntVar( &options.count, "n", 1, "Generate n combinations" )
    flag.StringVar( &options.outfile, "o", "-", "Specify an output file" )
    flag.BoolVar( &options.help, "h", false, "Print usage and default options" )
    flag.Parse()
}

func usage() {
    fmt.Fprintf(os.Stderr, "Usage : %s [OPTIONS] <production-file> \n", os.Args[0])
    flag.PrintDefaults()
}

func main() {
    arguments()
    options.prodfile = flag.Args()[0]
    if options.prodfile == ""  ||  options.help {
        usage()
        os.Exit(1)
    }

    options.random = rand.New( rand.NewSource( int64(options.seed) ))
    conf := make(golib.Config)
    start := monster.Parse(options.prodfile, conf)

    if options.ast {
        start.Show("")
    } else {
        c := make( monster.Context )
        nonterminals, root := monster.Build(start)
        c["_random"] = options.random
        c["_nonterminals"] = nonterminals
        for i:=0; i<options.count; i++ {
            outtext := root.Generate(c)
            fmt.Println(outtext)
        }
    }
}
