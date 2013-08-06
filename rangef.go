package monster
import ("fmt"; "strconv"; "math/rand")

func bnfrangef( c Context, nt *NonTerminal ) string {
    var min float64 = 0.0
    var max float64
    rnd := c["_random"].(*rand.Rand)
    cs := nt.Children // Arguments
    if len(cs) == 2 {
        min, _ = strconv.ParseFloat(cs[0].(*FloatTerminal).Value, 64)
        max, _ = strconv.ParseFloat(cs[1].(*FloatTerminal).Value, 64)
    } else if len(cs) == 1 {
        min, _ = strconv.ParseFloat(cs[0].(*FloatTerminal).Value, 64)
    } else {
        panic("Error: Atleast one argument expected in range() BNF")
    }
    return fmt.Sprintf( "%v", rnd.Float64() * (max-min) + min )
}

func init() {
    BnfCallbacks["rangef"] = bnfrangef
}

