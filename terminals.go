package monster

import (
	"fmt"
	"strconv"
)

// Show implements INode interface
func (t *Terminal) Show(prefix string) {
	fmt.Println(t.Repr(prefix))
}

// Repr implements INode interface
func (t *Terminal) Repr(prefix string) string {
	return fmt.Sprintf(prefix) + fmt.Sprintf("%v : %v ", t.Name, t.Value)
}

// Initialize implements INode interface
func (t *Terminal) Initialize(c Context) {
}

// Generate implements INode interface
func (t *Terminal) Generate(c Context) string {
	switch t.Name {
	case "STRING":
		return t.Value[1 : len(t.Value)-1]
	case "INT":
		val, _ := strconv.Atoi(t.Value)
		return fmt.Sprintf("%v", val)
	case "FLOAT":
		val, _ := strconv.ParseFloat(t.Value, 64)
		return fmt.Sprintf("%v", val)
	case "CHAR":
		return t.Value[1 : len(t.Value)-1]
	case "IDENT":
		return ""
	default:
		return t.Value
	}
}

// BNLTerminal represents a built-in-literal token
type BNLTerminal struct {
	Terminal
}

// Generate implements INode interface
func (t *BNLTerminal) Generate(c Context) string {
	return t.Value
}
