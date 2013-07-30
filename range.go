package monster

import "fmt"
import "os"
import "strconv"

func bnfrange( popts ParseOpts ) Terminal {
    var t Terminal
    if Token(popts.S) != "(" {
        fmt.Printf("bnf_range error")
        os.Exit(1)
    }
    min, max, toktype2 := "", "", ""
    toktype1, min := Tokent(popts.S)

    tok := Token(popts.S)
    if tok == "," {
        toktype2, max = Tokent(popts.S)
        tok = Token(popts.S)
    }
    if tok != ")" {
        fmt.Printf("Syntax error in bnf_range")
        os.Exit(1)
    }

    if toktype1 == "Float" || toktype2 == "Float" {
        minf, _ := strconv.ParseFloat(min,64)
        maxf, _ := strconv.ParseFloat(max,64)
        fn := func(context Context) string {
            return fmt.Sprintf( "%v", popts.Rnd.Float64() * (maxf-minf) + minf )
        }
        t = Terminal{ name: "BnfRange", value: "", generator: fn }
    } else {
        mini, _ := strconv.Atoi(min)
        maxi, _ := strconv.Atoi(max)
        fn := func(context Context) string {
            return fmt.Sprintf( "%v", popts.Rnd.Intn(maxi - mini) + mini )
        }
        t = Terminal{ name: "BnfRange", value: "", generator: fn }
    }
    return t
}

func init() {
    Bnfs["range"] = bnfrange
}
