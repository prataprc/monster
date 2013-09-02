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
    switch t.Name {
    case "String" : return t.Value[1: len(t.Value)-1]
    case "Int" : return fmt.Sprintf("%v", t.Value)
    case "Float" : return fmt.Sprintf("%v", t.Value)
    case "Char" : return t.Value[1: len(t.Value)-1]
    default : return t.Value
    }
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
