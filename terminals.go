package monster

// built-in Terminal constants
func nl() Terminal {   // Builtin terminal for double quotes
    return Terminal{
        name: "NL",
        value: `\n`,
        generator: func() string {return "\n"},
    }
}

func dq() Terminal {   // Builtin terminal for double quotes
    return Terminal{
        name: "DQ",
        value: `"`,
        generator: func() string {return `"`},
    }
}

func tRue() Terminal {   // Builtin terminal TRUE
    return Terminal{
        name: "TRUE",
        value: "true",
        generator: func() string {return "true"},
    }
}

func fAlse() Terminal {   // Builtin terminal FALSE
    return Terminal{
        name: "FALSE",
        value: "false",
        generator: func() string {return "false"},
    }
}

func null() Terminal {   // Builtin terminal NULL
    return Terminal{
        name: "NULL",
        value: "null",
        generator: func() string {return "null"},
    }
}

func init() {
    Terminals["NL"] = nl
    Terminals["DQ"] = dq
    Terminals["TRUE"] = tRue
    Terminals["FALSE"] = fAlse
    Terminals["NULL"] = null
}
