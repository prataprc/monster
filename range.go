package monster

import (
	"fmt"
	"math/rand"
	"strconv"
)

func bnfrange(c Context, args []interface{}) string {
	var min int = 0
	var max int
	rnd := c["_random"].(*rand.Rand)
	if len(args) == 2 {
		min, _ = strconv.Atoi(args[0].(string))
		max, _ = strconv.Atoi(args[1].(string))
	} else if len(args) == 1 {
		min, _ = strconv.Atoi(args[0].(string))
	} else {
		panic("Error: Atleast one argument expected in range() BNF")
	}
	return fmt.Sprintf("%v", rnd.Intn(max-min)+min)
}

func init() {
	BnfCallbacks["range"] = bnfrange
}
