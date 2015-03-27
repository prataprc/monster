//  Copyright (c) 2013 Couchbase, Inc.

package main

import "flag"
import "fmt"
import "log"
import "os"
import "time"
import "io/ioutil"

import "github.com/prataprc/goparsec"
import "github.com/prataprc/monster"
import "github.com/prataprc/monster/common"

var options struct {
	bagdir  string
	outfile string
	nonterm string
	seed    int
	count   int
	help    bool
	debug   bool
}

func argParse() (string, *os.File) {
	seed := time.Now().UTC().Second()
	flag.StringVar(&options.bagdir, "bagdir", "",
		"directory path containing bags")
	flag.StringVar(&options.nonterm, "nonterm", "",
		"evaluate the non-terminal")
	flag.IntVar(&options.seed, "seed", seed,
		"seed value")
	flag.IntVar(&options.count, "n", 1,
		"generate n combinations")
	flag.StringVar(&options.outfile, "o", "-",
		"specify an output file")
	flag.Parse()

	prodfile := flag.Args()[0]
	if prodfile == "" || options.help {
		usage()
		os.Exit(1)
	}

	var err error
	outfd := os.Stdout
	if options.outfile != "-" && options.outfile != "" {
		outfd, err = os.Create(options.outfile)
		if err != nil {
			log.Fatal(err)
		}
	}
	return prodfile, outfd
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage : %s [OPTIONS] <production-file> \n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	prodfile, outfd := argParse()

	// read production-file
	text, err := ioutil.ReadFile(prodfile)
	if err != nil {
		log.Fatal(err)
	}
	// compile
	root := compile(parsec.NewScanner(text))
	scope := root.(common.Scope)
	nterms := scope["_nonterminals"].(common.NTForms)
	if options.nonterm != "" {
		for i := 0; i < options.count; i++ {
			scope = monster.BuildContext(scope, uint64(i), options.bagdir)
			scope["_prodfile"] = prodfile
			val := evaluate("root", scope, nterms[options.nonterm])
			outtext := fmt.Sprintf("%v\n", val)
			if _, err := outfd.Write([]byte(outtext)); err != nil {
				log.Fatal(err)
			}
		}

	} else {
		// evaluate
		for i := 0; i < options.count; i++ {
			scope = monster.BuildContext(scope, uint64(i), options.bagdir)
			scope["_prodfile"] = prodfile
			val := evaluate("root", scope, nterms["s"])
			outtext := fmt.Sprintf("%v\n", val)
			if _, err := outfd.Write([]byte(outtext)); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func compile(s parsec.Scanner) parsec.ParsecNode {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%v at %v", r, s.GetCursor())
		}
	}()
	root, _ := monster.Y(s)
	return root
}

func evaluate(name string, scope common.Scope, forms []*common.Form) interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%v", r)
		}
	}()
	scope = scope.ApplyGlobalForms()
	return monster.EvalForms(name, scope, forms)
}
