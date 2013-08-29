package monster
import "strconv"
import "fmt"
import "reflect"
import "io/ioutil"
import "github.com/prataprc/golib/parsec"

type Interface interface{}
type Terminal struct {
    Name string         // typically contains terminal's token type
    Value string        // value of the terminal
    Tok parsec.Token    // Actual token obtained from the scanner
}
type NonTerminal struct {
    Name string         // typically contains terminal's token type
    Value string        // value of the terminal
    Children []INode
}
type INode interface{   // AST functions
    Show(string)
    Repr(prefix string) string
    Generate(c Context) string
}
type Context map[string]interface{}

var EMPTY = Terminal{Name: "EMPTY", Value:""}

var y = fmt.Sprintf("keep 'fmt' import during debugging"); // FIXME

//---- Global variables
// Built-in functions
var BnfCallbacks = make( map[string]func(Context, *NonTerminal)string )

func Parse(prodfile string, opts map[string]interface{}) INode {
    bs, _ := ioutil.ReadFile(prodfile)
    s := parsec.NewGoScan(bs, opts)
    return Y()(s).(INode)
}

func Build( start INode ) (map[string]INode, INode) {
    nonterminals := make(map[string]INode, 0)
    stcs := start.(*StartNT).Children
    for _, rb := range stcs {
        rbnt, _ := rb.(*RuleBlockNT)
        reflect.TypeOf(rbnt.Children[0]) // TODO : Remove this line
        // fmt.Println( reflect.TypeOf(rbnt.Children[0]) )
        nonterm, _ := rbnt.Children[0].(*IdentTerminal)
        if nonterm.Name != "Ident" {
            panic("Expected identifier !!")
        }
        nonterminals[ nonterm.Value ] = rbnt.Children[2]
    }
    root := stcs[0]
    return nonterminals, root
}

// parse nonterminal tokens
func Y() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            nt := NonTerminal{Name:"RULEBLOCKS", Children:inodes(ns)}
            return &StartNT{nt}
        }
        return parsec.Many( "y", nodify, true, ruleblock, parsec.End )()(s)
    }
}

func ruleblock() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            nt := NonTerminal{Name:"RULEBLOCK", Children:inodes(ns)}
            return &RuleBlockNT{nt}
        }
        return parsec.And(
            "ruleblock", nodify, true, ident, colon, rulelines, dot,
        )()(s)
    }
}

func rulelines() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            nt := NonTerminal{Name:"RULELINES", Children:inodes(ns)}
            return &RuleLinesNT{nt}
        }
        return parsec.Many( "rulelines", nodify, true, ruleline, pipe )()(s)
    }
}

func ruleline() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            nt := NonTerminal{Name:"RULELINE", Children:inodes(ns)}
            return &RuleLineNT{nt}
        }
        return parsec.And( "ruleline", nodify, true, rule, ruleOption )()(s)
    }
}

func rule() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            nt := NonTerminal{Name:"RULE", Children:inodes(ns)}
            return &RuleNT{nt}
        }
        nodifypart := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil {
                return nil
            } else {
                return ns[0]
            }
        }
        // dot, pipe and ruleOption Parsec should not be included.
        part := parsec.OrdChoice(
            "part", nodifypart, false, nl, dq, tRue, fAlse, null,
            literal, reference, bnf, ident,
        )
        return parsec.Many( "rule", nodify, true, part )()(s)
    }
}

func ruleOption() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodifyargs := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            return &NonTerminal{Name: "RULEARGS", Children:inodes(ns)}
        }
        nodify := func(n []parsec.ParsecNode) parsec.ParsecNode {
            var weight = 100
            cs := inodes(n)
            if len(cs) >= 3 {
                args := cs[1].(*NonTerminal)
                weight, _ = strconv.Atoi( args.Children[0].(*IntTerminal).Value )
            }
            nt := NonTerminal{Name: "RULEOPTIONS", Children:cs}
            return &RuleOptionsNT{nt, weight}
        }
        args := parsec.Many("ruleargs", nodifyargs, true, literal, comma)
        return parsec.And(
            "ruleoption", nodify, false, opencurl, args, closecurl,
        )()(s)
    }
}

func bnf() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodifyargs := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            return &NonTerminal{Name: "BNFARGS", Children:inodes(ns)}
        }
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            nt := NonTerminal{Name:"BNF", Children:inodes(ns)}
            return &BnfNT{nt}
        }
        bargs := parsec.Kleene( "bnfargs", nodifyargs, literal, comma )
        return parsec.And(
            "bnf", nodify, false, ident, openparan, bargs, closeparan,
        )()(s)
    }
}

func reference() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            nt := NonTerminal{Name:"REFERENCE", Children:inodes(ns)}
            return &ReferenceNT{nt}
        }
        return parsec.And( "reference", nodify, false, dollar, ident )()(s)
    }
}

// Terminal rats
func comment() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        tok := s.Peek(0)
        if tok.Type == "Comment" {
            s.Scan()
            return &CommentTerminal{terminal(tok)}
        } else {
            return nil
        }
    }
}

func ident() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        tok := s.Peek(0)
        if tok.Type == "Ident" {
            s.Scan()
            return &IdentTerminal{terminal(tok)}
        } else {
            return nil
        }
    }
}

var bang = prodTerminalize( "!", "BANG", "!" )
var hash = prodTerminalize( "#", "HASH", "#" )
var dot = prodTerminalize( ".", "DOT", "." )
var percent = prodTerminalize( "%", "PERCENT", "%" )
var colon = prodTerminalize( ":", "COLON", ":" )
var semicolon = prodTerminalize( ";", "SEMICOLON", ";" )
var comma = prodTerminalize( ",", "COMMA", "," )
var pipe = prodTerminalize( "|", "PIPE", "|" )
var dollar = prodTerminalize( "$", "DOLLAR", "$" )
var openparan = prodTerminalize( "(", "OPENPARAN", "(" )
var closeparan = prodTerminalize( ")", "CLOSEPARAN", ")" )
var opencurl = prodTerminalize( "{", "OPENPARAN", "{" )
var closecurl = prodTerminalize( "}", "CLOSEPARAN", "}" )

var nl = bnlTerminalize( "NL", "NEWLINE", "\n" )
var dq = bnlTerminalize( "DQ", "DQUOTE", "\"" )
var tRue = bnlTerminalize( "TRUE", "TRUE", "true" )
var fAlse = bnlTerminalize( "FALSE", "FALSE", "false" )
var null = bnlTerminalize( "NULL", "NULL", "null" )

// Parsec functions for terminals
func prodTerminalize(matchval string, n string, v string ) parsec.Parsec {
    return func() parsec.Parser {
        return func(s parsec.Scanner) parsec.ParsecNode {
            term := parsec.Terminalize(matchval, n, v)()(s)
            if term == nil {
                return nil
            } else {
                pt := term.(*parsec.Terminal)
                t := Terminal{Name:pt.Name, Value:pt.Value, Tok:pt.Tok}
                return &ProdTerminal{t}
            }
        }
    }
}

func bnlTerminalize(matchval string, n string, v string ) parsec.Parsec {
    return func() parsec.Parser {
        return func(s parsec.Scanner) parsec.ParsecNode {
            term := parsec.Terminalize(matchval, n, v)()(s)
            if term == nil {
                return nil
            } else {
                pt := term.(*parsec.Terminal)
                t := Terminal{Name:pt.Name, Value:pt.Value, Tok:pt.Tok}
                return &BNLTerminal{t}
            }
        }
    }
}

// Parsec functions to match `String`, `Char`, `Int`, `Float` literals
func literal() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        tok := s.Peek(0)
        if tok.Type == "String" {
            s.Scan()
            return &StrTerminal{terminal(tok)}
        } else if tok.Type == "Char" {
            s.Scan()
            return &CharTerminal{terminal(tok)}
        } else if tok.Type == "Int" {
            s.Scan()
            return &IntTerminal{terminal(tok)}
        } else if tok.Type == "Float" {
            s.Scan()
            return &FloatTerminal{terminal(tok)}
        } else {
            //fmt.Println(tok.Type, tok.Value)
            return nil
        }
    }
}

func inodes( pns []parsec.ParsecNode ) []INode {
    ins := make([]INode, 0)
    for _, n := range pns {
        ins = append( ins, n.(INode) )
    }
    return ins
}

func terminal(tok parsec.Token) Terminal {
    //fmt.Println(tok.Type, tok.Value)
    return Terminal{Name:tok.Type, Value:tok.Value, Tok:tok}
}
