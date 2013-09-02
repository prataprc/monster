package monster
import ("fmt"; "strconv"; "math/rand")

func bnfrange( c Context, nt *NonTerminal ) string {
    var min int = 0
    var max int
    rnd := c["_random"].(*rand.Rand)
    cs := nt.Children // Arguments
    if len(cs) == 2 {
        min, _ = strconv.Atoi(cs[0].(*Terminal).Value)
        max, _ = strconv.Atoi(cs[1].(*Terminal).Value)
    } else if len(cs) == 1 {
        min, _ = strconv.Atoi(cs[0].(*Terminal).Value)
    } else {
        panic("Error: Atleast one argument expected in range() BNF")
    }
    return fmt.Sprintf( "%v", rnd.Intn(max-min) + min )
}

func init() {
    BnfCallbacks["range"] = bnfrange
}
