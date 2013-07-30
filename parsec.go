package monster
import "text/scanner"

const (
    TERMINAL = atoi+1
    NONTERMINAL
    EMPTY
}
type Node struct {
    Type            // TERMINAL or NONTERMINAL or EMPTY
    T *Terminal
    NT *NonTerminal
}
type Terminal struct {
    Name string     // typically contains terminal's token type
    Value string    // value of the terminal
}
type NonTerminal struct {
    Name string
    Value interface{}
    Children []Node
}
type Token struct {
    Type string
    Value string
    pos scanner.Position
}

// This tokenizer is using text/scanner package. Make it generic so that
// parsec can be converted to a separate package.
func Tokenizer( sendto <-chan Token, src io.Reader ) {
    var s scanner.Scanner
    s.Init(src)
    tok := s.Scan()
    for tok != scanner.EOF {
        sendto <- Token {
            Type: scanner.TokenString(tok), Value: s.TokenText(), pos: s.Pos()
        }
        tok = s.Scan()
    }
    sendto <- Token {
        Type: scanner.TokenString(tok), Value: s.TokenText(), pos: s.Pos()
    }
}

func And( funs...interface{} ) {
}

func Or( funs...interface{} ) {
}

func Many( op, endOp ) {
}

func Kleene( op, endOp ) {
}

func Maybe( op ) {
}
