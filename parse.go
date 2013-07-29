package monster
import "text/scanner"
import "os"
import "strings"
import "fmt"
import "math/rand"

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
type ParseOpts struct {
    Prodfile string
    Rnd *rand.Rand
    S *scanner.Scanner
}

// Global variables
var Literals = make( map[string]func(string)Terminal )
var Terminals = make( map[string]func()Terminal )
var Bnfs = make( map[string]func(ParseOpts)Terminal )
var Bags = make( map[string][][]string )

func Parse(popts ParseOpts) (string, map[string]NonTerminal) {
    var s scanner.Scanner
    nonterminals := make( map[string]NonTerminal )
    file, err := os.Open(popts.Prodfile)
    if err != nil {
        fmt.Printf("Error with file")
    }
    popts.S = &s
    s.Init(file)
    ntname := Token( popts.S )
    start := ntname
    for i := 0; ntname != ""; i++ {
        rulename, nt := altrule(popts, ntname)
        nonterminals[ntname], ntname = nt, rulename
    }
    return start, buildAst(nonterminals)
}

func Token(s *scanner.Scanner) string {
    s.Scan()
    return s.TokenText()
}

func Tokent(s *scanner.Scanner) (string, string) {
    return scanner.TokenString(s.Scan()), s.TokenText()
}

func PrintAst( nonterminals map[string]NonTerminal ) {
    for ntname, nt := range nonterminals {
        fmt.Printf( "%v \n", ntname )
        for _, rule := range nt.any {
            printRule(rule)
        }
    }
}

func printRule( ns []Node ) {
    var outstr = make( []string, len(ns) )
    fmt.Printf("    : ")
    for i, n := range ns {
        outstr[i] = printNode(n)
    }
    fmt.Printf("%v\n", strings.Join(outstr, " . "))
    fmt.Printf("\n")
}

func printNode( n Node ) string {
    valt, ok := n.(Terminal)
    if ok {
        return fmt.Sprintf("%v(%v)", valt.name, valt.value)
    }
    valnt, ok := n.(NonTerminal)
    if ok {
        return fmt.Sprintf("%v", valnt.name)
    }
    fmt.Printf("Uknown Node %v", n)
    os.Exit(1)
    return ""
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

// Parse alternate rules for nonterminal `ntname`
func altrule( popts ParseOpts, ntname string ) (string, NonTerminal) {
    var any [100][]Node
    tok := Token(popts.S)
    i := 0
    for ; (tok == ":") || (tok == "|"); i++ {
        any[i] = nodes(popts)
        tok = Token(popts.S)
    }
    nt := NonTerminal{ name: ntname, value: "", generator: nil, any: any[:i] }
    return tok, nt
}

// Parse rule
func nodes( popts ParseOpts ) []Node {
    var all [100]Node
    var i int
    toktype, tok := Tokent(popts.S)
    for ; (tok != ".") && (toktype != "EOF") ; i++ {
        litfn := Literals[toktype]
        termfn := Terminals[tok]
        bnffn := Bnfs[tok]
        if litfn != nil {
            all[i] = litfn(tok)
        } else if termfn != nil {
            all[i] = termfn()
        } else if bnffn != nil {
            all[i] = bnffn( popts )
        } else {
            all[i] = tok
        }
        toktype, tok = Tokent(popts.S)
    }
    return all[:i]
}

// Build AST from map of `nonterminals`.
func buildAst( nonterminals map[string]NonTerminal ) map[string]NonTerminal {
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
    return nonterminals
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
