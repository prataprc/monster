* optimize for GC.
* fix this error,
    2015/09/25 17:09:21 listen tcp 127.0.0.1:6060: bind: address already in use
* Add definition in production grammar, to facilitate a generation on
  production rule with exact count. Eg,

  nonterminal : some "rule" {3}
              | alternate "rule".

  should generate `some "rule"` exactly 3 times before applying the
  alternation grammar for `nonterminal`.

* Remote the dependancy with golib.
* Add test cases.

