package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)
import "github.com/prataprc/monster"

var options struct {
	ast      bool
	prodfile string
	bagdir   string
	outfile  string
	seed     int
	random   *rand.Rand
	count    int
	help     bool
}

func argParse() {
	seed := time.Now().UTC().Second()
	flag.BoolVar(&options.ast, "ast", false,
		"show the ast of production")
	flag.StringVar(&options.bagdir, "bagdir", "",
		"directory path containing bags")
	flag.IntVar(&options.seed, "s", seed,
		"seed value")
	flag.IntVar(&options.count, "n", 1,
		"generate n combinations")
	flag.StringVar(&options.outfile, "o", "-",
		"specify an output file")
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage : %s [OPTIONS] <production-file> \n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	argParse()
	options.prodfile = flag.Args()[0]
	if options.prodfile == "" || options.help {
		usage()
		os.Exit(1)
	}

	options.random = rand.New(rand.NewSource(int64(options.seed)))
	conf := make(map[string]interface{})
	start, err := monster.Parse(options.prodfile, conf)
	if err != nil {
		panic(err)
	}

	if options.ast {
		start.Show("")
	} else {
		nonterminals, root := monster.Build(start)
		c := map[string]interface{}{
			"_random":       options.random,
			"_bagdir":       options.bagdir,
			"_nonterminals": nonterminals,
		}
		fd, err := os.Create(options.outfile)
		if err != nil {
			panic(err)
		}
		for i := 0; i < options.count; i++ {
			monster.Initialize(c)
			outtext := root.Generate(c) + "\n"
			if _, err := fd.Write([]byte(outtext)); err != nil {
				panic(err)
			} else {
				fmt.Println(outtext)
			}
		}
	}
}
