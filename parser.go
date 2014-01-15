package monster

import (
	"fmt"
	"github.com/prataprc/golib"
	"github.com/prataprc/goparsec"
	"io/ioutil"
	"strconv"
)

var _ = fmt.Sprintf("keep 'fmt' import during debugging") // FIXME

type Terminal struct {
	Name     string // typically contains terminal's token type
	Value    string // value of the terminal
	Position int
}

type NonTerminal struct {
	Name     string // typically contains terminal's token type
	Value    string // value of the terminal
	Children []parsec.ParsecNode
}

type INode interface { // AST functions
	Show(string)
	Repr(prefix string) string
	Initialize(c Context)
	Generate(c Context) string
}

type Context map[string]interface{}

var EMPTY = Terminal{Name: "EMPTY", Value: ""}

//---- Global variables
// Built-in functions
var BnfCallbacks = make(map[string]func(Context, []interface{}) string)

func Parse(prodfile string, conf golib.Config) INode {
	bs, _ := ioutil.ReadFile(prodfile)
	s := parsec.NewScanner(bs)
	root, _ := Y(s)
	return root.(INode)
}

func Build(start INode) (map[string]INode, INode) {
	nonterminals := make(map[string]INode)
	startnt := start.(*StartNT)
	root := startnt.Children[0].(INode)
	for _, nt := range startnt.Children {
		rb, _ := nt.(*RuleBlockNT)
		// fmt.Println( reflect.TypeOf(rb.Children[0]) )
		term := rb.Children[0].(*Terminal)
		nonterminals[term.Value] = rb.Children[1].(INode)
	}
	return nonterminals, root
}

func Initialize(c Context) {
	for _, node := range c["_nonterminals"].(map[string]INode) {
		node.Initialize(c)
	}
}

// Terminal rats
var bang = parsec.Token(`^\!`, "BANG")
var hash = parsec.Token(`^\#`, "HASH")
var dot = parsec.Token(`^\.`, "DOT")
var percent = parsec.Token(`^%`, "PERCENT")
var colon = parsec.Token(`^:`, "COLON")
var semicolon = parsec.Token(`^;`, "SEMICOLON")
var comma = parsec.Token(`^,`, "COMMA")
var pipe = parsec.Token(`^\|`, "PIPE")
var dollar = parsec.Token(`^\$`, "DOLLAR")
var openparan = parsec.Token(`^\(`, "OPENPARAN")
var closeparan = parsec.Token(`^\)`, "CLOSEPARAN")
var opencurl = parsec.Token(`^\{`, "OPENPARAN")
var closecurl = parsec.Token(`^\}`, "CLOSEPARAN")

var ident = parsec.Ident()

var nl = bnlToken("^NL", "NEWLINE", "\n")
var dq = bnlToken("^DQ", "DQUOTE", `"`)
var tRue = bnlToken("^TRUE", "TRUE", "true")
var fAlse = bnlToken("^FALSE", "FALSE", "false")
var null = bnlToken("^NULL", "NULL", "null")

// More terminal rats
func literal(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			t := ns[0].(*parsec.Terminal)
			return &Terminal{t.Name, t.Value, t.Position}
		}
		return nil
	}
	return parsec.OrdChoice(
		nodify, parsec.String(), parsec.Char(), parsec.Float(), parsec.Int(),
	)(s)
}

func bnlToken(matchval string, n string, v string) parsec.Parser {
	return func(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
		if term, news := parsec.Token(matchval, n)(s); term != nil {
			t := term.(*parsec.Terminal)
			term := Terminal{t.Name, v, t.Position}
			return &BNLTerminal{term}, news
		}
		return nil, s
	}
}

// nonterminal rats
func Y(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			return &StartNT{NonTerminal{Name: "RULEBLOCKS", Children: ns}}
		}
		return nil
	}
	return parsec.Many(nodify, ruleblock, parsec.NoEnd)(s)
}

func ruleblock(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			idt := ns[0].(*parsec.Terminal)
			idterm := Terminal{idt.Name, idt.Value, idt.Position}
			cs := []parsec.ParsecNode{&idterm, ns[2]}
			return &RuleBlockNT{NonTerminal{Name: "RULEBLOCK", Children: cs}}
		}
		return nil
	}
	return parsec.And(nodify, ident, colon, rulelines, dot)(s)
}

func rulelines(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			return &RuleLinesNT{NonTerminal{Name: "RULELINES", Children: ns}}
		}
		return nil
	}
	return parsec.Many(nodify, ruleline, pipe)(s)
}

func ruleline(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		var r *RuleNT
		var ropts *RuleOptionsNT
		if len(ns) > 0 {
			r = ns[0].(*RuleNT)
			if len(ns) > 1 {
				opts := ns[1].([]parsec.ParsecNode)
				if len(opts) > 0 {
					ropts = opts[0].(*RuleOptionsNT)
				}
			}
			return &RuleLineNT{NonTerminal{Name: "RULELINE"}, r, ropts}
		}
		return nil
	}
	return parsec.And(nodify, rule, parsec.Maybe(nil, ruleOption))(s)
}

func rule(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodifypart := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			if t, ok := ns[0].(*parsec.Terminal); ok {
				return &Terminal{t.Name, t.Value, t.Position}
			}
			return ns[0]
		}
		return nil
	}
	// Following order of parsers to OrdChoice is important, don't change !!
	part := parsec.OrdChoice(
		nodifypart, nl, dq, tRue, fAlse, null, literal, bnf, reference, ident,
	)

	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			return &RuleNT{NonTerminal{Name: "RULE", Children: ns}}
		}
		return nil
	}
	return parsec.Many(nodify, part)(s)
}

func ruleOption(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodifyargs := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			return &NonTerminal{Name: "RULEARGS", Children: ns}
		}
		return nil
	}
	args := parsec.Many(nodifyargs, literal, comma)

	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			sval := ns[1].(*NonTerminal).Children[0].(*Terminal).Value
			if weight, err := strconv.Atoi(sval); err == nil {
				return &RuleOptionsNT{weight, weight}
			} else {
				panic(err.Error())
			}
		}
		return nil

	}
	return parsec.And(nodify, opencurl, args, closecurl)(s)
}

func bnf(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodifyargs := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			return ns[0]
		}
		return nil
	}
	argparsers := parsec.OrdChoice(nodifyargs, literal, reference)
	bargs := parsec.Kleene(nil, argparsers, comma)

	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			return &BnfNT{
				NonTerminal{Name: "BNF"},
				ns[0].(*parsec.Terminal).Value, // CallName
				ns[2].([]parsec.ParsecNode),    // args
			}
		}
		return nil
	}
	return parsec.And(nodify, ident, openparan, bargs, closeparan)(s)
}

func reference(s parsec.Scanner) (parsec.ParsecNode, parsec.Scanner) {
	nodify := func(ns []parsec.ParsecNode) parsec.ParsecNode {
		if len(ns) > 0 {
			t := ns[1].(*parsec.Terminal)
			return &ReferenceNT{NonTerminal{Name: "REFERENCE", Value: t.Value}}
		}
		return nil
	}
	return parsec.And(nodify, dollar, ident)(s)
}
