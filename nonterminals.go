package monster
import ("strings"; "math/rand"; "fmt")

var x = fmt.Sprintf("keep 'fmt' import during debugging"); // FIXME

// Reference non-terminal
type ReferenceNT struct {
    NonTerminal
}
func (n *ReferenceNT) Generate(c Context) string {
    cs := n.Children
    return c[ cs[1].(*IdentTerminal).Value ].(string)
}

// built-in function non-terminal
type BnfNT struct {
    NonTerminal
}
func (n *BnfNT) Generate(c Context) string {
    cs := n.Children
    name := cs[0].(*IdentTerminal).Value
    return BnfCallbacks[ name ]( c, cs[2].(*NonTerminal) )
}

// rule-options non-terminal
type RuleOptionsNT struct {
    NonTerminal
    Weight int
}
func (n *RuleOptionsNT) Generate(c Context) string {
    return ""
}

// rule non-terminal
type RuleNT struct {
    NonTerminal
}
func (n *RuleNT) Generate(c Context) string {
    var s string
    keys := make( []string, len(n.Children) )[0:0]
    outs := make( []string, len(n.Children) )[0:0]
    for _, n := range n.Children {
        val, ok := n.(*IdentTerminal)
        if ok && val.Name == "Ident" {
            //fmt.Println(val.Value)
            n := c["_nonterminals"].(map[string]INode)[val.Value].(*RuleLinesNT)
            s = n.Generate(c)
            c[val.Value] = s
            keys = append( keys, val.Name )
        } else {
            s = n.Generate(c)
        }
        outs = append( outs, s )
    }
    for _, key := range keys {
        delete(c, key)
    }
    return strings.Join( outs, "" )
}

// rule-line non-terminal
type RuleLineNT struct {
    NonTerminal
}

// rule-lines non-terminal
type RuleLinesNT struct {
    NonTerminal
}
func (n *RuleLinesNT) Generate(c Context) string {
    ruleline := pickRuleLine( c, n.Children )
    return ruleline.Generate(c)
}
func pickRuleLine( c Context, cs []INode ) INode {
    var index = make(map[int]INode)
    accw := 0
    for _, ruleline := range cs {
        nt := ruleline.(*RuleLineNT).Children[1].(*RuleOptionsNT)
        if nt.Weight <= 0 { continue }
        accw += nt.Weight
        nt.Weight -= 1
        index[accw] = ruleline
    }
    if accw > 0 {
        r := c["_random"].(*rand.Rand).Intn(accw)
        for i, n := range index {
            if r <= i { return n }
        }
    }
    r := c["_random"].(*rand.Rand).Intn( len(cs) )
    return cs[r]
}


// rule-block non terminal
type RuleBlockNT struct {
    NonTerminal
}
func (n *RuleBlockNT) Generate(c Context) string {
    return n.Children[2].Generate(c)
}


// Start
type StartNT struct {
    NonTerminal
}
