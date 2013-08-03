package parsec
import ("testing"; "os";)
import "text/scanner"

func BenchmarkScanner(b *testing.B) {
    var s scanner.Scanner
    fd, _ := os.Open( "./sampleinput" )
    s.Init(fd)
    for {
        tok := Token {
                Type: scanner.TokenString( s.Scan() ),
                Value: s.TokenText(),
                Pos: s.Pos(),
        }
        if tok.Type == "EOF" { break }
    }
}

func BenchmarkGoscan(b *testing.B) {
    scanner := NewGoScanner( "./sampleinput" )
    for {
        tok := scanner.Scan()
        if tok.Type == "EOF" { break }
    }
}

