package monster

func str(value string) Terminal {   // string terminal
    value = value[1: len(value)-1] // remove the double quotes
    return Terminal{
        name: "String",
        value: value,
        generator: func() string {return value},
    }
}

func char(value string) Terminal {
    value = value[1: len(value)-1] // remove the single quotes
    return Terminal{
        name: "Char",
        value: value,
        generator: func() string {return value},
    }
}

func integer(value string) Terminal {
    return Terminal{
        name: "Int",
        value: value,
        generator: func() string {return value},
    }
}

func float(value string) Terminal {
    return Terminal{
        name: "Float",
        value: value,
        generator: func() string {return value},
    }
}

func init() {
    Literals["Int"] = integer
    Literals["String"] = str
    Literals["Char"] = char
    Literals["Float"] = float
}
