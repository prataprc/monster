package parsec
import "text/scanner"

type Scanner interface {
    Scan() Token
    Peek(int) Token
    Next() Token
    BookMark() int
    Rewind(int)
    Text() string
}

type Token struct {
    Type string
    Value string
    Pos scanner.Position
}

