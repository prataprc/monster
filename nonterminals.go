package monster

import (
	"fmt"
	"github.com/prataprc/goparsec"
	"math/rand"
	"strings"
)

var _ = fmt.Sprintf("keep 'fmt' import during debugging")

// INode interface for NonTerminal
func (n *NonTerminal) Show(prefix string) {
	for _, n := range n.Children {
		n.(*NonTerminal).Show(prefix + "  ")
	}
}
func (n *NonTerminal) Repr(prefix string) string {
	return fmt.Sprintf(prefix) + fmt.Sprintf("%v : %v \n", n.Name, n.Value)
}
func (n *NonTerminal) Initialize(c Context) {
	for _, child := range n.Children {
		if node, ok := child.(INode); ok {
			node.Initialize(c)
		} else {
			panic("Does not implement INode interface")
		}
	}
}
func (n *NonTerminal) Generate(c Context) string {
	s := ""
	for _, child := range n.Children {
		s += child.(*NonTerminal).Generate(c)
	}
	return s
}

// Start
type StartNT struct {
	NonTerminal
}

// rule-block non terminal
type RuleBlockNT struct {
	NonTerminal
}

func (n *RuleBlockNT) Generate(c Context) string {
	return n.Children[1].(*RuleLinesNT).Generate(c)
}

// rule-lines non-terminal
type RuleLinesNT struct {
	NonTerminal
}

func (n *RuleLinesNT) Generate(c Context) string {
	var index = make(map[int]*RuleLineNT)
	accw := 0
	for _, nt := range n.Children {
		ruleline := nt.(*RuleLineNT)
		ruleopts := ruleline.ruleopts
		if ruleopts != nil && ruleopts.weightCount > 0 {
			accw += ruleopts.weightCount
			ruleopts.weightCount -= 1
		} else {
			accw += 10 // Default weightage
		}
		index[accw] = ruleline
	}
	r := c["_random"].(*rand.Rand).Intn(accw)
	for i, ruleline := range index {
		if r < i {
			return ruleline.Generate(c)
		}
	}
	return n.Children[r].(*RuleLineNT).Generate(c)
}

// rule-line non-terminal
type RuleLineNT struct {
	NonTerminal
	rule     *RuleNT
	ruleopts *RuleOptionsNT
}

func (n *RuleLineNT) Initialize(c Context) {
	if n.rule != nil {
		n.rule.Initialize(c)
	}
	if n.ruleopts != nil {
		n.ruleopts.Initialize(c)
	}
}
func (n *RuleLineNT) Generate(c Context) string {
	return n.rule.Generate(c)
}

// rule non-terminal
type RuleNT struct {
	NonTerminal
}

func (n *RuleNT) Generate(c Context) string {
	var s string
	keys := make([]string, 0, len(n.Children))
	outs := make([]string, 0, len(n.Children))
	for _, child := range n.Children {
		if term, ok := child.(*Terminal); ok && (term.Name == "IDENT") {
			m := c["_nonterminals"].(map[string]INode)
			nonterm := m[term.Value].(*RuleLinesNT)
			s = nonterm.Generate(c)
			c[term.Value] = s
			keys = append(keys, term.Value)
		} else {
			s = child.(INode).Generate(c)
		}
		outs = append(outs, s)
	}
	for _, key := range keys {
		delete(c, key)
	}
	return strings.Join(outs, "")
}

// rule-options non-terminal
type RuleOptionsNT struct {
	weight      int
	weightCount int
}

func (n *RuleOptionsNT) Initialize(c Context) {
	n.weightCount = n.weight
}

// built-in function non-terminal
type BnfNT struct {
	NonTerminal
	callName string
	args     []parsec.ParsecNode
}

func (n *BnfNT) Generate(c Context) string {
	args := make([]interface{}, 0)
	for _, arg := range n.args {
		if term, ok := arg.(*Terminal); ok {
			args = append(args, term.Value)
		} else {
			refnt := arg.(*ReferenceNT)
			args = append(args, c[refnt.Value])
		}
	}
	return BnfCallbacks[n.callName](c, args)
}

// Reference non-terminal
type ReferenceNT struct {
	NonTerminal
}

func (n *ReferenceNT) Generate(c Context) string {
	return c[n.Value].(string)
}
