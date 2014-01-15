package monster

import (
	"fmt"
	"strconv"
)

// INode interface for Terminal
func (t *Terminal) Show(prefix string) {
	fmt.Println(t.Repr(prefix))
}
func (t *Terminal) Repr(prefix string) string {
	return fmt.Sprintf(prefix) + fmt.Sprintf("%v : %v ", t.Name, t.Value)
}
func (n *Terminal) Initialize(c Context) {
}
func (t *Terminal) Generate(c Context) string {
	switch t.Name {
	case "STRING":
		return t.Value[1 : len(t.Value)-1]
	case "INT":
		val, _ := strconv.Atoi(t.Value)
		return fmt.Sprintf("%v", val)
	case "FLOAT":
		val, _ := strconv.ParseFloat(t.Value, 64)
		fmt.Println(val, t.Value)
		return fmt.Sprintf("%v", val)
	case "CHAR":
		return t.Value[1 : len(t.Value)-1]
	case "IDENT":
		return ""
	default:
		return t.Value
	}
}

// Built-in terminal
type BNLTerminal struct {
	Terminal
}

func (t *BNLTerminal) Generate(c Context) string {
	return t.Value
}
