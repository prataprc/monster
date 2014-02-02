package monster

import (
	"fmt"
	"github.com/prataprc/goparsec"
	"math/rand"
	"strings"
)

// Show implements INode interface
func (n *NonTerminal) Show(prefix string) {
	for _, child := range n.Children {
		child.(INode).Show(prefix + "  ")
	}
}

// Repr implements INode interface
func (n *NonTerminal) Repr(prefix string) string {
	return fmt.Sprintf(prefix) + fmt.Sprintf("%v : %v \n", n.Name, n.Value)
}

// Initialize implements INode interface
func (n *NonTerminal) Initialize(c Context) {
	for _, child := range n.Children {
		if node, ok := child.(INode); ok {
			node.Initialize(c)
		} else {
			panic("Does not implement INode interface")
		}
	}
}

// Generate implements INode interface
func (n *NonTerminal) Generate(c Context) string {
	s := ""
	for _, child := range n.Children {
		s += child.(*NonTerminal).Generate(c)
	}
	return s
}

// StartNT represents the root node.
type StartNT struct {
	NonTerminal
}

// RuleBlockNT represents a rule block defining a set of one or more rulelines
// for a non-terminal
type RuleBlockNT struct {
	NonTerminal
}

// Show implements INode interface
func (n *RuleBlockNT) Show(prefix string) {
	n.Children[0].(INode).Show(prefix)
	for _, child := range n.Children[1:] {
		child.(INode).Show(prefix + "  ")
	}
}

// RuleLinesNT represents a set of one or more rulelines
type RuleLinesNT struct {
	NonTerminal
}

// Show implements INode interface
func (n *RuleLinesNT) Show(prefix string) {
	for i := 0; i < len(n.Children); i++ {
		child := n.Children[i]
		child.(INode).Show(prefix + "  ")
		if i < len(n.Children)-1 {
			fmt.Println(prefix + "|")
		}
	}
}

// Generate implements INode interface
func (n *RuleLinesNT) Generate(c Context) string {
	var index = make(map[int]*RuleLineNT)
	accw := 0
	for _, nt := range n.Children {
		ruleline := nt.(*RuleLineNT)
		ruleopts := ruleline.ruleopts
		if ruleopts != nil && ruleopts.weightCount > 0 {
			accw += ruleopts.weightCount
			ruleopts.weightCount--
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

// RuleLineNT represents a collection of rule definition that forms rule-line
// for a non-terminal. Also contains rule-options.
type RuleLineNT struct {
	NonTerminal
	rule     *RuleNT
	ruleopts *RuleOptionsNT
}

// Show implements INode interface
func (n *RuleLineNT) Show(prefix string) {
	n.rule.Show(prefix + "  ")
}

// Initialize implements INode interface
func (n *RuleLineNT) Initialize(c Context) {
	if n.rule != nil {
		n.rule.Initialize(c)
	}
	if n.ruleopts != nil {
		n.ruleopts.Initialize(c)
	}
}

// Generate implements INode interface
func (n *RuleLineNT) Generate(c Context) string {
	return n.rule.Generate(c)
}

// RuleNT represents an individual rule definition in a rule-line
type RuleNT struct {
	NonTerminal
}

// Generate implements INode interface
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

// RuleOptionsNT represents a rule option that defines the weight of a
// rule-line
type RuleOptionsNT struct {
	weight      int
	weightCount int
}

// Initialize implements INode interface
func (n *RuleOptionsNT) Initialize(c Context) {
	n.weightCount = n.weight
}

// BnfNT represents a built-in function reference in rule.
type BnfNT struct {
	NonTerminal
	callName string
	args     []parsec.ParsecNode
}

// Generate implements INode interface
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

// ReferenceNT represents a variable reference from context.
type ReferenceNT struct {
	NonTerminal
}

// Generate implements INode interface
func (n *ReferenceNT) Generate(c Context) string {
	return c[n.Value].(string)
}
