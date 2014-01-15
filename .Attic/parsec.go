package monster

import (
	"fmt"
)

type Parsec func(Scanner) INode
type Nodify func(*NonTerminal) INode

type Terminal struct {
	Name  string // typically contains terminal's token type
	Value string // value of the terminal
	Tok   Token  // Actual token obtained from the scanner
}
type NonTerminal struct {
	Name     string // typically contains terminal's token type
	Value    string // value of the terminal
	Children []INode
}

type INode interface { // AST functions
	Show(string)
	Generate(Context) string
}

var EMPTY = Terminal{Name: "EMPTY", Value: ""}

func And(name string, callb Nodify, assert bool, scanners ...Parsec) Parsec {
	return func(s Scanner) INode {
		var ns = make([]INode, 0)
		//fmt.Println(name)
		bm := s.BookMark()
		for _, scanner := range scanners {
			n := scanner(s)
			if n == nil && assert {
				Error(s, fmt.Sprintf("Error: %v \n", callb))
			} else if n == nil {
				s.Rewind(bm)
				return callb(nil)
			}
			ns = append(ns, n)
		}
		return callb(&NonTerminal{Children: ns})
	}
}

func OrdChoice(name string, scanners ...Parsec) Parsec {
	return func(s Scanner) INode {
		var n INode
		//fmt.Println(name, len(scanners))
		for _, scanner := range scanners {
			bm := s.BookMark()
			n = scanner(s)
			if n != nil {
				return n
			}
			s.Rewind(bm)
		}
		return nil
	}
}

func Kleene(name string, callb Nodify, scanners ...Parsec) Parsec {
	var opScan, sepScan Parsec
	opScan = scanners[0]
	if len(scanners) >= 2 {
		sepScan = scanners[1]
	}
	return func(s Scanner) INode {
		var ns = make([]INode, 0)
		//fmt.Println(name)
		for {
			n := opScan(s)
			if n == nil {
				break
			}
			ns = append(ns, n)
			if sepScan != nil && sepScan(s) == nil {
				break
			}
		}
		return callb(&NonTerminal{Children: ns})
	}
}

func Many(name string, callb Nodify, assert bool, scanners ...Parsec) Parsec {
	var opScan, sepScan Parsec
	opScan = scanners[0]
	if len(scanners) >= 2 {
		sepScan = scanners[1]
	}
	return func(s Scanner) INode {
		var ns = make([]INode, 0)
		//fmt.Println(name)
		bm := s.BookMark()
		n := opScan(s)
		if n == nil && assert {
			Error(s, fmt.Sprintf("atleast one of %v \n", callb))
		} else if n == nil {
			s.Rewind(bm)
			return callb(nil)
		} else {
			for {
				ns = append(ns, n)
				if sepScan != nil && sepScan(s) == nil {
					break
				}
				n = opScan(s)
				if n == nil {
					break
				}
			}
			return callb(&NonTerminal{Children: ns})
		}
		return nil
	}
}

func Maybe(name string, scanner Parsec) Parsec {
	return func(s Scanner) INode {
		n := scanner(s)
		//fmt.Println(name)
		return n
	}
}

func Error(s Scanner, str string) {
	panic(fmt.Sprintf("%v before %v \n", str, s.Next().Pos))
}

// INode interface for Terminal
func (t *Terminal) Show(prefix string) {
	fmt.Printf(prefix)
	fmt.Printf("%v : %v \n", t.Name, t.Value)
}
func (t *Terminal) Generate(c Context) string {
	return t.Value
}

// INode interface for NonTerminal
func (t *NonTerminal) Show(prefix string) {
	fmt.Printf(prefix)
	fmt.Printf("%v : %v \n", t.Name, t.Value)
	for _, n := range t.Children {
		n.Show(prefix + "  ")
	}
}
func (n *NonTerminal) Generate(c Context) string {
	s := ""
	for _, n := range n.Children {
		s += n.Generate(c)
	}
	return s
}
