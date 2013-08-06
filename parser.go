package monster
import ("strconv"; "fmt"; "reflect")

type Context map[string]interface{}

var y = fmt.Sprintf("keep 'fmt' import during debugging"); // FIXME

//---- Global variables
// Built-in functions
var BnfCallbacks = make( map[string]func(Context, *NonTerminal)string )

func Parse(prodfile string, parser Parsec) INode {
    s := NewGoScan(prodfile)
    return parser(s)
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
func Y() Parsec {
    nodify := func(n *NonTerminal) INode {
        return &StartNT{NonTerminal{Name:"RULEBLOCKS", Children:n.Children}}
    }
    return Many("y", nodify, true, ruleblock(), end())
}

func ruleblock() Parsec {
    nodify := func(n *NonTerminal) INode {
        // fmt.Printf("Parsing nonterminal %v ... \n", n.Children[0].(*IdentTerminal).Value)
        return &RuleBlockNT{NonTerminal{Name:"RULEBLOCK", Children:n.Children}}
    }
    return And( "ruleblock", nodify, true, ident(), colon(), rulelines(), dot())
}

func rulelines() Parsec {
    nodify := func(n *NonTerminal) INode {
        return &RuleLinesNT{NonTerminal{Name:"RULELINES", Children:n.Children}}
    }
    return Many("rulelines", nodify, true, ruleline(), pipe())
}

func ruleline() Parsec {
    nodify := func(n *NonTerminal) INode {
        return &RuleLineNT{NonTerminal{Name:"RULELINE", Children:n.Children}}
    }
    return And( "ruleline", nodify, true, rule(), ruleOption() )
}

func rule() Parsec {
    nodify := func(n *NonTerminal) INode {
        return &RuleNT{NonTerminal{Name:"RULE", Children:n.Children}}
    }
    // dot, pipe and ruleOption Parsec should not be included.
    part := OrdChoice(
                "part", nl(), dq(), tRue(), fAlse(), null(), reference(), bnf(),
                ident(), literal() )
    return Many( "rule", nodify, true, part )
}

func ruleOption() Parsec {
    nodifyargs := func(n *NonTerminal) INode {
        n.Name = "RULEARGS"
        return n
    }
    nodify := func(n *NonTerminal) INode {
        var cs []INode
        var weight = 100
        if n != nil {
            cs = n.Children
        }
        if len(cs) >= 3 {
            args := cs[1].(*NonTerminal)
            weight, _ = strconv.Atoi( args.Children[0].(*IntTerminal).Value )
        }
        return &RuleOptionsNT{NonTerminal{Name:"RULEOPTIONS"},weight}
    }
    args := Many( "ruleargs", nodifyargs, true, literal(), comma() )
    return And( "ruleoption", nodify, false, opencurl(), args, closecurl() )
}

func bnf() Parsec {
    nodifyargs := func(n *NonTerminal) INode {
        n.Name = "BNFARGS"
        return n
    }
    nodify := func(n *NonTerminal) INode {
        if n == nil { return nil }
        return &BnfNT{NonTerminal{Name:"BNF", Children:n.Children}}
    }
    bargs := Kleene( "bnfargs", nodifyargs, literal(), comma())
    return And( "bnf", nodify, false, ident(), openparan(), bargs, closeparan())
}

func reference() Parsec {
    nodify := func(n *NonTerminal) INode {
        if n == nil { return nil }
        return &ReferenceNT{NonTerminal{Name:"REFERENCE", Children:n.Children}}
    }
    return And( "reference", nodify, false, dollar(), ident() )
}

// Terminal rats
func comment() Parsec {
    return func(s Scanner) INode {
        tok := s.Peek(0)
        //fmt.Println("comment")
        if tok.Type == "Comment" {
            s.Scan()
            t := Terminal{Name:tok.Type, Value:tok.Value, Tok:tok}
            return &CommentTerminal{t}
        } else {
            return nil
        }
    }
}

func ident() Parsec {
    return func(s Scanner) INode {
        tok := s.Peek(0)
        //fmt.Println("ident")
        if tok.Type == "Ident" {
            s.Scan()
            t := Terminal{Name: tok.Type, Value: tok.Value, Tok: tok}
            return &IdentTerminal{t}
        } else {
            return nil
        }
    }
}

func literal() Parsec {
    return func(s Scanner) INode {
        //fmt.Println("literal")
        tok := s.Peek(0)
        t := Terminal{Name: tok.Type, Value: tok.Value, Tok: tok}
        if tok.Type == "String" {
            s.Scan()
            return &StrTerminal{t}
        } else if tok.Type == "Char" {
            s.Scan()
            return &CharTerminal{t}
        } else if tok.Type == "Int" {
            s.Scan()
            return &IntTerminal{t}
        } else if tok.Type == "Float" {
            s.Scan()
            return &FloatTerminal{t}
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
func prodTerminalize(matchval string, n string, v string ) func()Parsec {
    return func() Parsec {
        return func(s Scanner) INode {
            //fmt.Println("prodterm")
            tok := s.Peek(0)
            if matchval == tok.Value {
                s.Scan()
                return &ProdTerminal{Terminal{Name: n, Value: v, Tok: tok}}
            } else {
                return nil
            }
        }
    }
}

func bnlTerminalize(matchval string, n string, v string ) func()Parsec {
    return func() Parsec {
        return func(s Scanner) INode {
            //fmt.Println("bnlterm")
            tok := s.Peek(0)
            if matchval == tok.Value {
                s.Scan()
                return &BNLTerminal{Terminal{Name: n, Value: v, Tok: tok }}
            } else {
                return nil
            }
        }
    }
}

func end() Parsec {
    return func(s Scanner) INode {
        tok := s.Next()
        if tok.Type == "EOF" {
            return nil
        }
        return &Terminal{Name: tok.Type, Value: tok.Value, Tok: tok }
    }
}
