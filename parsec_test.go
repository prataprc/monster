package monster
import ("testing"; "os"; "io/ioutil")
import "text/scanner"

var testfile = "./sampletest"

func BenchmarkPrepare(b *testing.B) {
    data, err := ioutil.ReadFile("./parsec_test.go")
    fd, err := os.Create( testfile )
    defer func(){ fd.Close() }()
    if err == nil {
        for i:=0; i<10000; i++ {
            fd.Write(data)
        }
    }
}

func BenchmarkScanner(b *testing.B) {
    var s scanner.Scanner
    fd, _ := os.Open( testfile )
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
    scanner := NewGoScan(testfile)
    for {
        tok := scanner.Scan()
        if tok.Type == "EOF" { break }
    }
}

func BenchmarkUnprepare(b *testing.B) {
    os.Remove( testfile )
}

