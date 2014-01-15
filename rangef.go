package monster

import (
	"fmt"
	"math/rand"
	"strconv"
)

func bnfrangef(c Context, args []interface{}) string {
	var min float64 = 0.0
	var max float64
	rnd := c["_random"].(*rand.Rand)
	if len(args) == 2 {
		min, _ = strconv.ParseFloat(args[0].(string), 64)
		max, _ = strconv.ParseFloat(args[1].(string), 64)
	} else if len(args) == 1 {
		min, _ = strconv.ParseFloat(args[0].(string), 64)
	} else {
		panic("Error: Atleast one argument expected in range() BNF")
	}
	return fmt.Sprintf("%v", rnd.Float64()*(max-min)+min)
}

func init() {
	BnfCallbacks["rangef"] = bnfrangef
}
