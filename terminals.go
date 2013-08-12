package monster
import "fmt"

// INode interface for Terminal
func (t *Terminal) Show( prefix string ) {
    fmt.Println( t.Repr(prefix) )
}
func (t *Terminal) Repr( prefix string ) string {
    return fmt.Sprintf(prefix) + fmt.Sprintf("%v : %v ", t.Name, t.Value)
}
func (t *Terminal) Generate(c Context) string {
    return t.Value
}

// Indentifier terminal
type CommentTerminal struct {
    Terminal
}
func (t *CommentTerminal) Generate(c Context) string {
    return ""
}

// Indentifier terminal
type IdentTerminal struct {
    Terminal
}
func (t *IdentTerminal) Generate(c Context) string {
    return ""
}

// String terminal
type StrTerminal struct {
    Terminal
}
func (t *StrTerminal) Generate(c Context) string {
    value := t.Value
    return value[1: len(value)-1] // remove the double quotes
}

// Integer terminal
type IntTerminal struct {
    Terminal
}
func (t *IntTerminal) Generate(c Context) string {
    return fmt.Sprintf("%v", t.Value)
}

// Float terminal
type FloatTerminal struct {
    Terminal
}
func (t *FloatTerminal) Generate(c Context) string {
    return fmt.Sprintf("%v", t.Value)
}

// Character terminal
type CharTerminal struct {
    Terminal
}
func (t *CharTerminal) Generate(c Context) string {
    value := t.Value
    return value[1: len(value)-1] // remove the single quotes
}

// Built-in terminal
type BNLTerminal struct {
    Terminal
}

// Production terminal
type ProdTerminal struct {
    Terminal
}
func (t *ProdTerminal) Generate(c Context) string {
    return ""
}
