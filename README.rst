Monster is a grammar based production tool. At present I am using this to
generate json documents using production grammar.

Can be invoked via command-line like,

.. code-block:: bash

    $ go run monster/monster.go -h
        -ast=false: show the ast of production
        -bagdir="": directory path containing bags
        -n=1: generate n combinations
        -o="-": specify an output file
        -s=15: seed value

Example production grammar, that generate randomized user documents.

.. code-block:: bnf

    json : "{ " properties "}".
    properties  : DQ "type" DQ ": " DQ "user" DQ ", " NL
                  DQ "first-name" DQ ": " DQ fname DQ ", " NL
                  lastname
                  age
                  emailid
                  city
                  gender.

    fname       : bag("./propernames").
    lastname    : DQ "last-name"  DQ ": " DQ bag("./propernames") DQ ", " NL.
    age         : DQ "age"        DQ ": " range(15, 80)                     ", " NL.
    emailid     : DQ "emailid"    DQ ": " DQ $fname "@gmail.com"     DQ ", " NL.
    city        : DQ "city"       DQ ": " DQ bag("./cities") DQ ", " NL.
    gender      : DQ "gender"     DQ ": " DQ "male" DQ NL
                | DQ "gender"     DQ ": " DQ "female" DQ NL.


The grammar can be invoked via command line like,

.. code-block:: bash

    go run monster/monster.go -bagdir bags/ prods/json.prod

use `-n` switch to generate as many document as needed, documents will be output
to stdout by default, to redirect them to a file use `-o` switch.
