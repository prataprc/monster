package parsec
import "text/scanner"
import "io"
import "fmt"
import "os"

type Goscan struct {
    text string
    req chan<- interface{}
    res <-chan interface{}
    filename string
}

func NewGoScan(filename string) *Goscan {
    var res = make( chan interface{} )
    var req  = make( chan interface{} )
    var fd *os.File
    fd, _ = os.Open(filename)
    text = fd.read()
    fd.close()
    fd, _ = os.Open(filename)
    go doscan( req, res, fd )
    return &Goscan{ req: req, res: res, filename:filename }
}

func (s *Goscan) Text() string {
    return s.text
}

func (s *Goscan) Scan() Token {
    var cmd = make([]interface{}, 1)
    cmd[0] = "scan"
    s.req<-cmd
    return (<-s.res).(Token)
}

func (s *Goscan) Next() Token {
    var cmd = make([]interface{}, 1)
    cmd[0] = "next"
    s.req<- cmd
    return (<-s.res).(Token)
}

func (s *Goscan) Peek(offset int) Token {
    var cmd = make([]interface{}, 2)
    cmd[0] = "peek"
    cmd[1] = offset
    s.req <- cmd
    return (<-s.res).(Token)
}

func (s *Goscan) BookMark() int {
    var cmd = make([]interface{}, 2)
    cmd[0] = "bookmark"
    s.req <- cmd
    return (<-s.res).(int)
}

func (s *Goscan) Rewind(offset int) {
    var cmd = make([]interface{}, 2)
    cmd[0] = "rewind"
    cmd[1] = offset
    s.req <- cmd
    <-s.res
}


// This tokenizer is using text/scanner package. Make it generic so that
// parsec can be converted to a separate package.
func doscan( req <-chan interface{}, res chan<- interface{}, src io.Reader ) {
    var s scanner.Scanner
    var curtok = 0

    s.Init(src)
    toks := fullscan(&s)
    for {
        cmd := (<-req).([]interface{})
        switch cmd[0].(string) {
        case "bookmark" :
            res <- curtok
        case "rewind" :
            curtok = cmd[1].(int)
            res <- Token{} // Dummy
        case "scan" :
            res <- toks[curtok]
            curtok += 1
        case "peek" :
            off := cmd[1].(int)
            if off >= 1 { panic("Offset to peek should be greater than 0") }
            res <- toks[curtok+off]
        case "next" :
            res <- toks[curtok+1]
        default :
            fmt.Printf("Unknown command to goscan : %v\n", cmd[0].(string))
        }
    }
}

func fullscan( s *scanner.Scanner ) []Token {
    var toks = make( []Token, 0 )
    for {
        tok := Token {
            Type: scanner.TokenString( s.Scan() ),
            Value: s.TokenText(),
            Pos: s.Pos(),
        }
        if tok.Type == "EOF" { break }
        toks = append(toks, tok )
    }
    return toks
}
