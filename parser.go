//  Copyright (c) 2013 Couchbase, Inc.

// Parser uses `parsec` tool to parse production grammar and
// construct generator tree.

package monster

import "fmt"
import "log"
import "time"
import "strconv"
import "math/rand"

import "github.com/prataprc/goparsec"
import "github.com/prataprc/monster/builtin"
import "github.com/prataprc/monster/common"

//----------------
// Monster grammar
//----------------

// bnf        -> forms nterminal*
// nterminal  -> ident ":" rules "."
// rules      -> ruletok+
//            -> rules "|" ruletok+
// ruletok    -> ident
//            |  ref
//            |  terminal
//            |  string
//            |  form
//
// forms      -> form*
// form       -> "(" formarg+ ")"

// ident      -> `[a-z0-9]+`
// terminal   -> `[A-Z][A-Z0-9]*`
// formarg    -> `[^ \t\r\n\(\)]+`
//            |  form
// ws         -> `[ \t\r\n]+`

// Nt is intermediate data structure.
type Nt [2]interface{}

// EvalForms refers to common.EvalForms
var EvalForms = common.EvalForms

// Circular rats
var form parsec.Parser

// Y root combinator for monster.
var Y parsec.Parser

// Terminal rats
var formtok = parsec.Token(`[^ \t\r\n\(\)]+`, "FORMTOK")
var ident = parsec.Token(`[a-z0-9]+`, "IDENT")
var ref = parsec.Token(`[$#][a-z0-9]+`, "REF")
var term = parsec.Token(`[A-Z][A-Z0-9]*`, "TERM")
var sTring = parsec.String()
var literaltok = parsec.OrdChoice(
	litNode,
	parsec.Float(), parsec.Hex(), parsec.Oct(), parsec.Int(),
	parsec.String(),
	parsec.Token(`true`, "TRUE"), parsec.Token(`false`, "FALSE"))
var openparan = parsec.Token(`\(`, "OPENPARAN")
var closeparan = parsec.Token(`\)`, "CLOSEPARAN")
var dot = parsec.Token(`\.`, "DOT")
var colon = parsec.Token(`\:`, "COLON")
var pipe = parsec.Token(`\|`, "PIPE")

// NonTerminal rats
var formarg = parsec.OrdChoice(formtokNode, literaltok, ref, formtok, &form)
var ruletok = parsec.OrdChoice(ruletokNode, ident, term, sTring, ref, &form)
var rules = parsec.Many(rulesNode, parsec.Many(ruleNode, ruletok, nil), pipe)
var nterm = parsec.And(ntermNode, ident, colon, rules, dot)

func init() {
	form = parsec.And(
		formNode,
		openparan, ident, parsec.Kleene(nil, formarg, nil), closeparan)
	Y = parsec.And(rootNode,
		parsec.Kleene(formsNode, &form, nil),
		parsec.Kleene(ntermsNode, nterm, nil))

	initBuiltins()
	initLiterals()
}

// BuildContext to initialize a new scope for evaluating production grammars.
func BuildContext(
	scope common.Scope,
	seed uint64,
	bagdir string) common.Scope {

	rootns := []parsec.ParsecNode{scope["_globalForms"], scope["_nonterminals"]}
	scope = common.NewScopeFromRoot(rootns)
	scope.SetBagdir(bagdir)
	if seed != 0 {
		scope.SetRandom(rand.New(rand.NewSource(int64(seed))))
	} else {
		now := time.Now().UnixNano()
		scope.SetRandom(rand.New(rand.NewSource(int64(now))))
	}
	for _, name := range scope.FormDuplicates(builtins) {
		log.Printf("warning: `%v` non-terminal is defined as builtin\n", name)
	}
	return scope
}

func rootNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	return common.NewScopeFromRoot(ns)
}

func formsNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	formls := make([]*common.Form, 0, len(ns))
	for _, n := range ns {
		formls = append(formls, n.(*common.Form))
	}
	return formls
}

func formNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	name := ns[1].(*parsec.Terminal).Value
	ns = ns[2].([]parsec.ParsecNode)
	form, ok := builtins[name]
	if ok { // apply builtin form.
		return common.NewForm(
			name,
			func(scope common.Scope, _ ...interface{}) interface{} {
				args := make([]interface{}, 0, len(ns))
				for _, n := range ns {
					args = append(args, n.(*common.Form).Eval(scope))
				}
				return form.Eval(scope, args...)
			})
	}
	// apply non-terminal
	return common.NewForm(
		"#"+name,
		func(scope common.Scope, _ ...interface{}) interface{} {
			forms, ok := scope.GetNonTerminal(name)
			if ok {
				val := EvalForms(name, scope, forms)
				scope.Set(name, val, false /*global*/)
				return val
			}
			panic(fmt.Errorf("unknown form name %v\n", name))
		})
}

func ntermsNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	ntls := make(common.NTForms)
	for _, n := range ns {
		v := n.(Nt)
		ntls[v[0].(string)] = v[1].([]*common.Form)
	}
	return ntls
}

func ntermNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	t := ns[0].(*parsec.Terminal)
	return Nt{t.Value, ns[2]}
}

func rulesNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	rulels := make([]*common.Form, 0, len(ns))
	weight := 1.0 / float64(len(ns))
	for _, n := range ns {
		rule := n.(*common.Form)
		rule.SetDefaultWeight(weight)
		rulels = append(rulels, rule)
	}
	return rulels
}

func ruleNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	// compute rule weight.
	var weight, restrain float64
	if len(ns) > 1 {
		if weigh := ns[0].(*common.Form); weigh.Name == "weigh" {
			rs := weigh.Eval(make(common.Scope)).([]interface{})
			weight, restrain = rs[0].(float64), rs[1].(float64)
			ns = ns[1:]
		}
	}
	// compose rule-form.
	rats := make([]*common.Form, 0, len(ns))
	for _, n := range ns {
		rats = append(rats, n.(*common.Form))
	}
	form := common.NewForm(
		"##rule",
		func(scope common.Scope, _ ...interface{}) interface{} {
			str := ""
			for i, rat := range rats {
				val := rat.Eval(scope)
				if val == nil {
					return nil
				}
				scope.Set("#"+strconv.Itoa(i), val, false /*global*/)
				str += fmt.Sprintf("%v", val)
			}
			return str
		})
	form.SetWeight(weight, restrain)
	return form
}

func ruletokNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	switch n := ns[0].(type) {
	case *parsec.Terminal:
		switch n.Name {
		case "IDENT":
			return identNode(n)
		case "TERM":
			return termNode(n)
		case "REF":
			return refNode(n)
		}

	case string:
		str := n[1 : len(n)-1]
		return common.NewForm(
			"##string",
			func(_ common.Scope, _ ...interface{}) interface{} { return str })

	case *common.Form:
		return n
	}
	panic(fmt.Errorf("unknown form type %T\n", ns[0]))
}

func formtokNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	switch n := ns[0].(type) {
	case *parsec.Terminal:
		switch n.Name {
		case "TERM":
			return termNode(n)
		case "REF":
			return refNode(n)
		case "FORMTOK":
			return common.NewForm(
				"##formtok",
				func(_ common.Scope, _ ...interface{}) interface{} {
					return n.Value
				})
		}

	case string:
		str := n[1 : len(n)-1]
		return common.NewForm(
			"##string",
			func(_ common.Scope, _ ...interface{}) interface{} { return str })

	case *common.Form:
		return n
	}
	panic(fmt.Errorf("unknown form type %T\n", ns[0]))
}

func varNode(n *parsec.Terminal) *common.Form {
	return common.NewForm(
		"##var",
		func(scope common.Scope, _ ...interface{}) interface{} {
			val, _, ok := scope.Get(n.Value)
			if !ok {
				panic(fmt.Errorf("unknown variable %v\n", n.Value))
			}
			return val
		})
}

func refNode(n *parsec.Terminal) *common.Form {
	return common.NewForm(
		"##ref",
		func(scope common.Scope, _ ...interface{}) interface{} {
			switch n.Value[0] {
			case '$':
				val, _, ok := scope.Get(n.Value[1:])
				if !ok {
					panic(fmt.Errorf("unknown reference %v\n", n.Value))
				}
				return val
			case '#':
				val, _, ok := scope.Get(n.Value)
				if !ok {
					panic(fmt.Errorf("unknown argument %v\n", n.Value))
				}
				return val
			}
			panic(fmt.Errorf("unknown form %v as part of rule\n", n.Value))
		})
}

func identNode(n *parsec.Terminal) *common.Form {
	return common.NewForm(
		"##ident",
		func(scope common.Scope, _ ...interface{}) interface{} {
			name := n.Value
			forms, ok := scope.GetNonTerminal(name)
			if ok {
				val := EvalForms(name, scope, forms)
				scope.Set(n.Value, val, false /*global*/)
				return val
			}
			panic(fmt.Errorf("unknown nonterminal %v\n", n.Value))
		})
}

func termNode(n *parsec.Terminal) *common.Form {
	str, _ := literals[n.Value]
	return common.NewForm(
		"##term",
		func(_ common.Scope, _ ...interface{}) interface{} { return str })
}

func litNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) == 0 {
		return nil
	}

	if s, ok := ns[0].(string); ok {
		return s
	}

	var val interface{}
	var err error
	t := ns[0].(*parsec.Terminal)
	switch t.Name {
	case "INT":
		val, err = strconv.ParseInt(t.Value, 10, 64)
		if err != nil {
			panic(fmt.Errorf("cannot parse %v for integer\n", t))
		}

	case "HEX":
		val, err = strconv.ParseInt(t.Value, 16, 64)
		if err != nil {
			panic(fmt.Errorf("cannot parse %v for hexadecimal\n", t))
		}

	case "OCT":
		val, err = strconv.ParseInt(t.Value, 8, 64)
		if err != nil {
			panic(fmt.Errorf("cannot parse %v for octal\n", t))
		}

	case "FLOAT":
		val, err = strconv.ParseFloat(t.Value, 64)
		if err != nil {
			panic(fmt.Errorf("cannot parse %v for float64\n", t))
		}

	case "STRING":
		val = t.Value[1 : len(t.Value)-1]

	case "TRUE":
		val = true

	case "FALSE":
		val = false
	}
	return common.NewForm(
		"##literaltok",
		func(_ common.Scope, _ ...interface{}) interface{} { return val })
}

func one2one(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) == 0 {
		return nil
	}
	return ns[0]
}

//--------------------
// initialize builtins
//--------------------

var builtins = make(map[string]*common.Form)
var literals = make(map[string]string)

func initBuiltins() {
	builtins["let"] = common.NewForm("let", builtin.Let)
	builtins["letr"] = common.NewForm("letr", builtin.Letr)
	builtins["global"] = common.NewForm("global", builtin.Global)
	builtins["weigh"] = common.NewForm("weigh", builtin.Weigh)
	builtins["bag"] = common.NewForm("bag", builtin.Bag)
	builtins["range"] = common.NewForm("range", builtin.Range)
	builtins["rangef"] = common.NewForm("rangef", builtin.Rangef)
	builtins["ranget"] = common.NewForm("ranget", builtin.Ranget)
	builtins["choice"] = common.NewForm("choice", builtin.Choice)
	builtins["uuid"] = common.NewForm("uuid", builtin.Uuid)
	builtins["inc"] = common.NewForm("inc", builtin.Inc)
	builtins["dec"] = common.NewForm("dec", builtin.Dec)
	builtins["len"] = common.NewForm("len", builtin.Len)
	builtins["sprintf"] = common.NewForm("sprintf", builtin.Sprintf)
}

func initLiterals() {
	literals["DQ"] = "\""
	literals["NL"] = "\n"
}
