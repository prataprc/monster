//  Copyright (c) 2013 Couchbase, Inc.

package main

import "bytes"
import "encoding/json"
import "flag"
import "fmt"
import "log"
import "os"
import "time"
import "unsafe"
import "math/rand"
import "strings"
import "reflect"
import "io/ioutil"
import _ "net/http/pprof"
import "net/http"
import "runtime/pprof"
import "runtime/debug"

import "github.com/prataprc/goparsec"
import "github.com/prataprc/monster"
import "github.com/prataprc/monster/common"

var options struct {
	bagdir   string
	outfile  string
	nonterms []string
	mprof    string
	pprof    string
	seed     int
	count    int
	par      int
	help     bool
	json     bool
	debug    bool
}

func argParse() (string, *os.File) {
	var nonterms string

	seed := time.Now().UTC().Second()
	flag.StringVar(&options.bagdir, "bagdir", "",
		"directory path containing bags")
	flag.StringVar(&nonterms, "nonterms", "s",
		"comma seperated list of non-terminals to pick one of them as root")
	flag.StringVar(&options.mprof, "mprof", "",
		"dump mem-profile to file")
	flag.StringVar(&options.pprof, "pprof", "",
		"dump cpu-profile to file")
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
	options.nonterms = strings.Split(nonterms, ",")

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
	go func() { log.Println(http.ListenAndServe("localhost:6061", nil)) }()

	// read production-file
	text, err := ioutil.ReadFile(prodfile)
	if err != nil {
		log.Fatal(err)
	}

	if options.pprof != "" {
		fd, err := os.Create(options.pprof)
		if err != nil {
			log.Fatalf("unable to create %q: %v\n", options.pprof, err)
		}
		defer fd.Close()

		pprof.StartCPUProfile(fd)
		defer pprof.StopCPUProfile()
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
	if takeMEMProfile(options.mprof) {
		fmt.Printf("dumped mem-profile to %v\n", options.mprof)
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
	rnd := rand.New(rand.NewSource(int64(seed)))
	for i := 0; i < count; i++ {
		nonterm := options.nonterms[rnd.Intn(len(options.nonterms))]
		scope = scope.RebuildContext()
		val := []byte(evaluate("root", scope, nterms[nonterm]).(string))
		if !options.json {
			outch <- val
		} else if err := json.Unmarshal(val, &value); err == nil {
			outch <- val
		} else {
			log.Fatalf("Invalid JSON %v\n", err)
		}
	}
}

func compile(s parsec.Scanner) parsec.ParsecNode {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%v at %v\n", r, s.GetCursor())
			fmt.Printf("%v\n", getStackTrace(2, debug.Stack()))
		}
	}()
	root, _ := monster.Y(s)
	return root
}

func evaluate(name string, scope common.Scope, forms []*common.Form) interface{} {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%v", r)
			fmt.Printf("%v\n", getStackTrace(2, debug.Stack()))
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

func str2bytes(str string) []byte {
	if str == "" {
		return nil
	}
	st := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sl := &reflect.SliceHeader{Data: st.Data, Len: st.Len, Cap: st.Len}
	return *(*[]byte)(unsafe.Pointer(sl))
}

func getStackTrace(skip int, stack []byte) string {
	var buf bytes.Buffer
	lines := strings.Split(string(stack), "\n")
	for _, call := range lines[skip*2:] {
		buf.WriteString(fmt.Sprintf("%s\n", call))
	}
	return buf.String()
}
