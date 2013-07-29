package monster
import "text/scanner"
import "strconv"
import "os"
import "strings"
import "fmt"
import "math/rand"
import "encoding/csv"

type Node interface{}
// Types of Terminal as indicated by `name` field,
//  BnfRange, BnfBag
type Terminal struct {
    name string
    value string
    generator func() string
}
type NonTerminal struct {
    name string
    value string
    generator func() string
    any [][]Node
}

// Global variables
type ParseOpts struct {
    Prodfile string
    Rnd *rand.Rand
    S *scanner.Scanner
}
var bags = make( map[string][][]string )

func Parse(popts ParseOpts) (string, map[string]NonTerminal) {
    var s scanner.Scanner
    nonterminals := make( map[string]NonTerminal )
    file, err := os.Open(popts.Prodfile)
    if err != nil {
        fmt.Printf("Error with file")
    }
    popts.S = &s
    s.Init(file)
    ntname := token( popts.S )
    for i := 0; ntname != ""; i++ {
        rulename, nt := altrule(popts, ntname)
        nonterminals[ntname], ntname = nt, rulename
    }
    return buildAst( nonterminals )
}

func PrintAst( nonterminals map[string]NonTerminal ) {
    for ntname, nt := range nonterminals {
        fmt.Printf( "%v \n", ntname )
        for _, rule := range nt.any {
            printRule(rule)
        }
    }
}

func Generate( popts ParseOpts, nt NonTerminal ) string {
    rule := nt.any[ popts.Rnd.Intn(len(nt.any)) ] // Randomly pick a rule
    terms := make( []string, len(rule) )
    for _, node := range rule  {
        val, ok := node.(Terminal)
        if ok { // Node is terminal
            // fmt.Printf("-- terminal `%v`\n", val.name)
            terms = append(terms, val.generator())
        } else { // Node is nonterminal
            val, _ := node.(NonTerminal)
            // fmt.Printf("-- nonterminal `%v`\n", val.name)
            terms = append( terms, Generate(popts, val) )
        }
    }
    return strings.Join( terms, "" )
}

func printRule( ns []Node ) {
    fmt.Printf( "    : ")
    for _, n := range ns { printNode( n ) }
    fmt.Printf("\n")
}

func printNode( n Node ) {
    valt, ok := n.(Terminal)
    if ok {
        fmt.Printf("%v(%v) . ", valt.name, valt.value)
        return
    }
    valnt, ok := n.(NonTerminal)
    if ok {
        fmt.Printf("%v . ", valnt.name)
        return
    }
    fmt.Printf("Uknown Node %v", n)
    os.Exit(1)
}

// Local functions
func altrule( popts ParseOpts, ntname string ) (string, NonTerminal) {
    var any [100][]Node
    tok := token(popts.S)
    i := 0
    for ; (tok == ":") || (tok == "|"); i++ {
        any[i] = nodes(popts)
        tok = token(popts.S)
    }
    nt := NonTerminal{ name: ntname, value: "", generator: nil, any: any[:i] }
    return tok, nt
}

func nodes( popts ParseOpts ) []Node {
    var all [100]Node
    var i int
    toktype, tok := tokent(popts.S)
    for ; (tok != ".") && (toktype != "EOF") ; i++ {
        if toktype == "String" {
            all[i] = t_string(tok)
        } else if toktype == "Int" {
            all[i] = t_int(tok)
        } else if toktype == "Float" {
            all[i] = t_float(tok)
        } else if toktype == "Char" {
            all[i] = t_char(tok)
        } else if tok == "NL" {
            all[i] = t_NL()
        } else if tok == "DQ" {
            all[i] = t_DQ()
        } else if tok == "TRUE" {
            all[i] = t_TRUE()
        } else if tok == "FALSE" {
            all[i] = t_FALSE()
        } else if tok == "NULL" {
            all[i] = t_NULL()
        } else {
            t, ok := bnf( popts, tok )
            if ok {
                all[i] = t
            } else {
                all[i] = tok
            }
        }
        toktype, tok = tokent(popts.S)
    }
    return all[:i]
}

func bnf( popts ParseOpts, value string ) (Terminal, bool) {
    switch value {
    case "bag" : return bnf_bag(popts), true
    case "range" : return bnf_range(popts), true
    }
    return Terminal{}, false
}

func bnf_bag( popts ParseOpts ) Terminal {
    var filename string
    var index int = 0
    var t Terminal
    _ = token(popts.S) // "("
    toktype, filename := tokent(popts.S)
    if toktype != "String" {
        fmt.Printf("Error: bag() first argument should be string")
        os.Exit(1)
    }
    filename = filename[1: len(filename)-1] // remove the double quotes
    tok := token(popts.S)
    if tok == "," {
        toktype, tok = tokent(popts.S)
        if toktype != "Int" {
            fmt.Printf("Error: bag() second argument should be integer")
            os.Exit(1)
        }
        index, _ = strconv.Atoi( tok )
        tok = token(popts.S)
    }
    if tok == ")" {
        fn := func() string { return rangeOnFile(popts, filename, index) }
        return Terminal{ name : "BnfBag", value : "", generator: fn }
    } else {
        fmt.Printf("Syntax error in bnf_bag")
        os.Exit(1)
    }
    return t // Dummy return, otherwise compiler cribs
}

func bnf_range( popts ParseOpts ) Terminal {
    var t Terminal
    if token(popts.S) != "(" {
        fmt.Printf("bnf_range error")
        os.Exit(1)
    }
    min, max, toktype2 := "", "", ""
    toktype1, min := tokent(popts.S)

    tok := token(popts.S)
    if tok == "," {
        toktype2, max = tokent(popts.S)
        tok = token(popts.S)
    }
    if tok != ")" {
        fmt.Printf("Syntax error in bnf_range")
        os.Exit(1)
    }

    if toktype1 == "Float" || toktype2 == "Float" {
        minf, _ := strconv.ParseFloat(min,64)
        maxf, _ := strconv.ParseFloat(max,64)
        fn := func() string {
            return fmt.Sprintf( "%v", popts.Rnd.Float64() * (maxf-minf) + minf )
        }
        t = Terminal{ name: "BnfRange", value: "", generator: fn }
    } else {
        mini, _ := strconv.Atoi(min)
        maxi, _ := strconv.Atoi(max)
        fn := func() string { 
            return fmt.Sprintf( "%v", popts.Rnd.Intn(maxi - mini) + mini )
        }
        t = Terminal{ name: "BnfRange", value: "", generator: fn }
    }
    return t
}

// literals as Terminals
func t_string(value string) Terminal {   // string terminal
    value = value[1: len(value)-1] // remove the double quotes
    return Terminal{ name: "String", value: value, generator: func() string {return value} }
}

func t_char(value string) Terminal {
    value = value[1: len(value)-1] // remove the single quotes
    return Terminal{ name: "Char", value: value, generator: func() string {return value} }
}

func t_int(value string) Terminal {
    return Terminal{ name: "Int", value: value, generator: func() string {return value} }
}

func t_float(value string) Terminal {
    return Terminal{ name: "Float", value: value, generator: func() string {return value} }
}

// built-in Terminal constants
func t_NL() Terminal {   // Builtin terminal for double quotes
    return Terminal{ name: "NL", value: `\n`, generator: func() string {return "\n"} }
}

func t_DQ() Terminal {   // Builtin terminal for double quotes
    return Terminal{ name: "DQ", value: `"`, generator: func() string {return `"`} }
}

func t_TRUE() Terminal {   // Builtin terminal TRUE
    return Terminal{ name: "TRUE", value: "true", generator: func() string {return "true"} }
}

func t_FALSE() Terminal {   // Builtin terminal FALSE
    return Terminal{ name: "FALSE", value: "false", generator: func() string {return "false"} }
}

func t_NULL() Terminal {   // Builtin terminal NULL
    return Terminal{ name: "NULL", value: "null", generator: func() string {return "null"} }
}

// local helper functions
func rangeOnFile(popts ParseOpts, filename string, index int) string {
    var choice = bags[filename]
    if choice == nil {
        choice = readBag( filename )
        bags[filename] = choice
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

func buildAst( nonterminals map[string]NonTerminal ) (string, map[string]NonTerminal) {
    keys := nt( nonterminals )
    nons := make([]string, len(keys)-1)
    for _, nt := range nonterminals {
        for i, rule := range nt.any {
            for j, node := range rule {
                val, ok := node.(string)
                if ok {
                    //fmt.Printf("^`%v` `%v` `%v`\n", nt.name, val, nonterminals[val].name)
                    nt.any[i][j] = nonterminals[val]
                    nons = append(nons, val)
                }
            }
        }
    }
    start := ""
    for _, key := range keys {
        if contains(key, nons) { continue }
        if start != "" {
            fmt.Printf( "More than one Start element %v %v\n", start, key )
            os.Exit(1)
        }
        start = key
    }
    return start, nonterminals
}

func token(s *scanner.Scanner) string {
    s.Scan()
    return s.TokenText()
}

func tokent(s *scanner.Scanner) (string, string) {
    return scanner.TokenString(s.Scan()), s.TokenText()
}


// Not yet used.
func colBag( records [][]string, index int ) []string {
    values, i := make([]string, len(records)), 0
    for _, record := range records {
        values[i] = record[index]
        i++
    }
    return values
}

func contains( key string, arrs []string ) bool {
    for _, v := range arrs {
        if key == v { return true }
    }
    return false
}

func nt( m map[string]NonTerminal ) []string {
    var keys = make([]string, len(m))
    i := 0
    for key := range m {
        keys[i] = key
        i++
    }
    return keys
}
