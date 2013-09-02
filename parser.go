package monster
import "strconv"
import "fmt"
import "reflect"
import "io/ioutil"
import "github.com/prataprc/golib/parsec"
import "github.com/prataprc/golib"

type Terminal struct {
    parsec.Terminal
}
type NonTerminal struct {
    parsec.NonTerminal
    Children []INode
}
type Interface interface{}
type INode interface{   // AST functions
    Show(string)
    Repr(prefix string) string
    Generate(c Context) string
}
type Context map[string]interface{}

var EMPTY = Terminal{parsec.Terminal{Name: "EMPTY", Value:""}}

var _ = fmt.Sprintf("keep 'fmt' import during debugging"); // FIXME

//---- Global variables
// Built-in functions
var BnfCallbacks = make( map[string]func(Context, *NonTerminal)string )

func Parse(prodfile string, conf golib.Config) INode {
    bs, _ := ioutil.ReadFile(prodfile)
    s := parsec.NewGoScan(bs, conf)
    return Y()(s).(INode)
}

func Build(start INode) (map[string]INode, INode) {
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
            return &StartNT{
                NonTerminal{parsec.NonTerminal{Name:"RULEBLOCKS"}, inodes(ns)},
            }
        }
        return parsec.Many( "y", nodify, true, ruleblock, parsec.NoEnd )()(s)
    }
}

func ruleblock() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            return &RuleBlockNT{
                NonTerminal{parsec.NonTerminal{Name:"RULEBLOCK"}, inodes(ns)},
            }
        }
        return parsec.And(
            "ruleblock", nodify, true, ident, colon, rulelines, dot,
        )()(s)
    }
}

func rulelines() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            return &RuleLinesNT{
                NonTerminal{parsec.NonTerminal{Name:"RULELINES"}, inodes(ns)},
            }
        }
        return parsec.Many( "rulelines", nodify, true, ruleline, pipe )()(s)
    }
}

func ruleline() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            return &RuleLineNT{
                NonTerminal{parsec.NonTerminal{Name:"RULELINE"}, inodes(ns)},
            }
        }
        return parsec.And( "ruleline", nodify, true, rule, ruleOption )()(s)
    }
}

func rule() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            return &RuleNT{
                NonTerminal{parsec.NonTerminal{Name:"RULE"}, inodes(ns)},
            }
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
            return &NonTerminal{
                parsec.NonTerminal{Name:"RULEARGS"}, inodes(ns),
            }
        }
        nodify := func(n []parsec.ParsecNode) parsec.ParsecNode {
            var weight = 100
            cs := inodes(n)
            if len(cs) >= 3 {
                args := cs[1].(*NonTerminal)
                weight, _ = strconv.Atoi( args.Children[0].(*Terminal).Value )
            }
            return &RuleOptionsNT{
                NonTerminal{parsec.NonTerminal{Name:"RULEOPTIONS"}, cs},
                weight,
            }
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
            return &NonTerminal{
                parsec.NonTerminal{Name:"BNFARGS"}, inodes(ns),
            }
        }
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns == nil { return nil }
            return &BnfNT{
                NonTerminal{parsec.NonTerminal{Name:"BNF"}, inodes(ns)},
            }
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
            return &ReferenceNT{
                NonTerminal{parsec.NonTerminal{Name:"REFERENCE"}, inodes(ns)},
            }
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

func literal() parsec.Parser {
    return func(s parsec.Scanner) parsec.ParsecNode {
        nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
            if ns != nil {
                return &Terminal{*ns[0].(*parsec.Terminal)}
            } else {
                return nil
            }
        }
        return parsec.OrdChoice(
            "literal", nodify, false,
            parsec.Literalof("String"), parsec.Literalof("Char"),
            parsec.Literalof("Int"), parsec.Literalof("Float"),
        )()(s)
    }
}

func prodToken(matchval string, n string) parsec.Parsec {
    return func() parsec.Parser {
        return func(s parsec.Scanner) parsec.ParsecNode {
            if term := parsec.Tokenof(matchval, n)()(s); term !=nil {
                return &ProdTerminal{Terminal{*term.(*parsec.Terminal)}}
            }
            return nil
        }
    }
}

func bnlToken(matchval string, n string, v string) parsec.Parsec {
    return func() parsec.Parser {
        return func(s parsec.Scanner) parsec.ParsecNode {
            if term := parsec.Tokenof(matchval, n)()(s); term != nil {
                t := term.(*parsec.Terminal)
                t.Value = v
                return &BNLTerminal{Terminal{*t}}
            }
            return nil
        }
    }
}

var bang = prodToken( `\!`, "BANG" )
var hash = prodToken( `\#`, "HASH" )
var dot = prodToken( `\.`, "DOT" )
var percent = prodToken( `%`, "PERCENT" )
var colon = prodToken( `:`, "COLON" )
var semicolon = prodToken( `;`, "SEMICOLON" )
var comma = prodToken( `,`, "COMMA" )
var pipe = prodToken( `\|`, "PIPE" )
var dollar = prodToken( `\$`, "DOLLAR" )
var openparan = prodToken( `\(`, "OPENPARAN" )
var closeparan = prodToken( `\)`, "CLOSEPARAN" )
var opencurl = prodToken( `\{`, "OPENPARAN" )
var closecurl = prodToken( `\}`, "CLOSEPARAN" )

var nl = bnlToken( "NL", "NEWLINE", "\n" )
var dq = bnlToken( "DQ", "DQUOTE", "\"" )
var tRue = bnlToken( "TRUE", "TRUE", "true" )
var fAlse = bnlToken( "FALSE", "FALSE", "false" )
var null = bnlToken( "NULL", "NULL", "null" )

func inodes( pns []parsec.ParsecNode ) []INode {
    ins := make([]INode, 0)
    for _, n := range pns {
        // fmt.Println(n)
        ins = append( ins, n.(INode) )
    }
    return ins
}

func terminal(tok *parsec.Token) Terminal {
    //fmt.Println(tok.Type, tok.Value)
    return Terminal{parsec.Terminal{Name:tok.Type, Value:tok.Value, Tok:tok}}
}
