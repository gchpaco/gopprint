# gopprint
An implementation of Kiselyov et al's pretty printing algorithm in Go.

Pretty printing occurs in two global phases.  Rather than try to
print some random tree directly, which could get quite ugly
quite quickly, we build a "pretty printer document" out of some
very simple primitives.  These primitives (and our algorithm) are
due to Oppen[1] originally and later Kiselyov[2].  Oppen's original
formulation had `Text`, `LineBreak`, `Concat`, and `Group`.  I
generalized `LineBreak` to `cond_t` which became our `cond_t`
because we need to do more sophisticated breaks, and I added
`nest_t` for controllable indentation.
[1]: Oppen, D.C.: Prettyprinting. ACM Trans. Program. Lang. Syst. 2
     (1980) 465–483.  Not available online without an ACM subscription.

[2]: Kiselyov, O., Petyon-Jones, S. and Sabry, A.: Lazy v. Yield:
     Incremental, Linear Pretty-printing.  Available online at
     http://okmij.org/ftp/continuations/PPYield/yield-pp.pdf
