package parsec

type Node interface{}
type Parsec func(Scanner) Node
type Terminal struct {
    Name string     // typically contains terminal's token type
    Value string    // value of the terminal
    Tok Token       // Actual token obtained from the scanner
}
type NonTerminal struct {
    Name string
    Value interface{}
    Children []Node
}
var EMPTY = Terminal{Name: "EMPTY", Value:""}

func And( name string, ops ...Parsec ) Parsec {
    return func( s Scanner ) Node {
        var ns = make([]Node, 0)
        bm := s.BookMark()
        for _, op := range ops {
            n := op(s)
            if n == nil {
                s.Rewind(bm)
                return nil
            }
            ns = append(ns, n)
        }
        return NonTerminal{Name:name, Children:ns}
    }
}

func OrdChoice( ops ...Parsec ) Parsec {
    return func(s Scanner) Node {
        for _, op := range ops {
            bm := s.BookMark()
            n := op(s)
            if n != nil {
                return n
            }
            s.Rewind(bm)
        }
        return nil
    }
}

func Kleene( name string, ops ...Parsec ) Parsec {
    return func(s Scanner) Node {
        var ns = make([]Node, 0)
        op, sepOp := ops[0], ops[1]
        for {
            n := op(s)
            if n == nil {
                break
            }
            ns = append(ns, n)
            if sepOp != nil  &&  sepOp(s) == nil {
                break
            }
        }
        return NonTerminal{Name:name, Children: ns}
    }
}

func Many( name string, ops ...Parsec) Parsec {
    return func(s Scanner) Node {
        var ns = make([]Node, 0)
        bm := s.BookMark()
        op, sepOp := ops[0], ops[1]
        n := op(s)
        if n == nil {
            s.Rewind(bm)
            return nil
        } else {
            for {
                ns = append(ns, n)
                if sepOp != nil  &&  sepOp(s) == nil {
                    break
                }
                n := op(s)
                if n == nil {
                    break
                }
            }
            return NonTerminal{Name: name, Children: ns}
        }
    }
}

func Maybe( name string, op Parsec ) Parsec {
    return func(s Scanner) Node {
        n := op(s)
        if n == nil {
            return Terminal{Name:"EMPTY", Value:""}
        } else {
            return n
        }
    }
}
