Monster is a grammar based production tool. At present I am using this to
generate json documents using production grammar.

Can be invoked via command-line like,

.. code-block:: bash

    $ go run monster/monster.go -h
          -bagdir="": directory path containing bags
          -n=1: generate n combinations
          -nonterm="": evaluate the non-terminal
          -o="-": specify an output file
          -seed=37: seed value

Example production grammar, that generate randomized user documents.

.. code-block:: bnf

    s : value.

    object : "{" properties "}".

    properties : properties "," property
               | property.

    property   : DQ (bag "./web2") DQ ":" value.

    array   : "[" values "]".

    values  : value "," value
            | (weigh 0.8 0.2) values.

    value   : (weigh 0.1) basic
            | (weigh 0.45 0.1) array
            | (weigh 0.45 0.1) object.

    basic   : "true"
            | "false"
            | "null"
            | number
            | string.

    string  : DQ (bag "./web2") DQ.

    number  : (range 0 100000)
            | (rangef 0.0 100.0).

The grammar can be invoked via command line like,

.. code-block:: bash

    go run monster/monster.go -bagdir ./bags -n 10 ./prods/json.prod

use `-n` switch to generate as many document as needed, documents will be output
to stdout by default, to redirect them to a file use `-o` switch.
