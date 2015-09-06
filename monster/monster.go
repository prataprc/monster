//  Copyright (c) 2013 Couchbase, Inc.

package main

import "encoding/json"
import "flag"
import "fmt"
import "log"
import "os"
import "time"
import "io/ioutil"
import _ "net/http/pprof"
import "net/http"
import "runtime/pprof"

import "github.com/prataprc/goparsec"
import "github.com/prataprc/monster"
import "github.com/prataprc/monster/common"

var options struct {
	bagdir  string
	outfile string
	nonterm string
	memprof string
	seed    int
	count   int
	par     int
	help    bool
	json    bool
	debug   bool
}

func argParse() (string, *os.File) {
	seed := time.Now().UTC().Second()
	flag.StringVar(&options.bagdir, "bagdir", "",
		"directory path containing bags")
	flag.StringVar(&options.nonterm, "nonterm", "s",
		"evaluate the non-terminal")
	flag.StringVar(&options.memprof, "memprof", "",
		"dump mem-profile to file")
	flag.IntVar(&options.seed, "seed", seed,
		"seed value")
	flag.IntVar(&options.count, "count", 1,
		"generate count number of combinations")
	flag.IntVar(&options.par, "par", 1,
		"number of parallel routines to use to generate")
	flag.StringVar(&options.outfile, "o", "-",
		"specify an output file")
	flag.BoolVar(&options.json, "json", false,
		"type of data to generate - json or string")
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

	// start pprof
	go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()

	// read production-file
	text, err := ioutil.ReadFile(prodfile)
	if err != nil {
		log.Fatal(err)
	}

	// spawn generators
	outch := make(chan []byte, 1000)
	for i := 0; i < options.par; i++ {
		go generate(text, options.count, prodfile, outch)
	}

	// gather generated values
	till := options.count * options.par
	for till > 0 {
		outtext := <-outch
		outtext = append(outtext, '\n')
		if _, err := outfd.Write([]byte(outtext)); err != nil {
			log.Fatal(err)
		}
		till--
	}
	fmt.Println("Completed !!")
	if takeMEMProfile(options.memprof) {
		fmt.Printf("dumped mem-profile to %v\n", options.memprof)
	}
}

func generate(text []byte, count int, prodfile string, outch chan<- []byte) {
	// compile
	root := compile(parsec.NewScanner(text)).(common.Scope)
	seed, bagdir, prodfile := uint64(options.seed), options.bagdir, prodfile
	scope := monster.BuildContext(root, seed, bagdir, prodfile)
	nterms := scope["_nonterminals"].(common.NTForms)

	// verify the sanity of json generated from production file
	var value map[string]interface{}
	if options.json {
		scope = scope.RebuildContext()
		val := evaluate("root", scope, nterms[options.nonterm])
		if err := json.Unmarshal([]byte(val.(string)), &value); err != nil {
			log.Printf("Invalid JSON %v\n", err)
			os.Exit(1)
		} else {
			outch <- []byte(val.(string))
		}
	}

	for i := 0; i < count; i++ {
		scope = scope.RebuildContext()
		val := evaluate("root", scope, nterms[options.nonterm])
		outch <- []byte(val.(string))
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
	return monster.EvalForms(name, scope, forms)
}

func takeMEMProfile(filename string) bool {
	if filename == "" {
		return false
	}
	fd, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
		return false
	}
	pprof.WriteHeapProfile(fd)
	defer fd.Close()
	return true
}
