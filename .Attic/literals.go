package monster
import "fmt"
import "parsec"

func str(value string) parsec.Terminal {   // string terminal
    value = value[1: len(value)-1] // remove the double quotes
    return parsec.Terminal{ Name: "String", Value: value }
}

func char(value string) parsec.Terminal {
    value = value[1: len(value)-1] // remove the single quotes
    return parsec.Terminal{ Name: "Char", Value: value }
}

func integer(value string) parsec.Terminal {
    return Terminal{ Name: "Int", Value: value }
}

func float(value string) parsec.Terminal {
    return Terminal{ Name: "Float", Value: value }
}

func reference(value string) parsec.Terminal {
    fn := func(context Context) string {
        return fmt.Sprintf("%v", context[value])
    }
    return Terminal {
        name : "Reference",
        value : value,
        generator: fn,
    }
}

func init() {
    Literals["Int"] = integer
    Literals["String"] = str
    Literals["Char"] = char
    Literals["Float"] = float
    Literals["Reference"] = reference
}
