package monster
import (
    "strings"
    "math/rand"
    "fmt"
    "github.com/prataprc/golib/parsec"
)

var _ = fmt.Sprintf("keep 'fmt' import during debugging");

// INode interface for NonTerminal
func (n *NonTerminal) Show(prefix string) {
    fmt.Printf( "%v", n.Repr(prefix) )
    for _, n := range n.Children {
        n.(*NonTerminal).Show(prefix + "  ")
    }
}
func (n *NonTerminal) Repr( prefix string ) string {
    return fmt.Sprintf(prefix) + fmt.Sprintf("%v : %v \n", n.Name, n.Value)
}
func (n *NonTerminal) Generate(c Context) string {
    s := ""
    for _, n := range n.Children {
        s += n.(*NonTerminal).Generate(c)
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
        if ruleopts != nil && ruleopts.weight > 0 {
            accw += ruleopts.weight
            ruleopts.weight -= 1
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
    rule *RuleNT
    ruleopts *RuleOptionsNT
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
            keys = append(keys, term.Name)
        } else {
            s = child.(INode).Generate(c)
        }
        outs = append(outs, s)
    }
    for _, key := range keys {
        delete(c, key)
    }
    return strings.Join( outs, "" )
}

// rule-options non-terminal
type RuleOptionsNT struct {
    weight int
}

// built-in function non-terminal
type BnfNT struct {
    NonTerminal
    callName string
    args []parsec.ParsecNode
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

