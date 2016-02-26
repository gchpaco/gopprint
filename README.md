# gopprint
An implementation of Kiselyov et al's pretty printing algorithm in Go.

Pretty printing occurs in two global phases.  Rather than try to print
some random tree directly, which could get quite ugly quite quickly,
we build a "pretty printer document" out of some very simple
primitives.  These primitives (and our algorithm) are due to
D.C. Oppen originally and later Kiselyov et al.  Oppen's original
formulation had `Text`, `LineBreak`, `Concat`, and `Group`.  I
generalized `LineBreak` to `cond_t` which became our `cond_t` because
we need to do more sophisticated breaks, and I added `nest_t` for
controllable indentation.

 * http://dl.acm.org.sci-hub.io/citation.cfm?id=357115 Oppen, D.C.: Prettyprinting. ACM Trans. Program. Lang. Syst. 2 (1980) 465–483.    Not available online without an ACM subscription.
 * http://okmij.org/ftp/continuations/PPYield/yield-pp.pdf Kiselyov, O., Peyton-Jones, S. and Sabry, A.: Lazy v. Yield: Incremental, Linear Pretty-printing.
