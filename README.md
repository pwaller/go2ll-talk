go2ll-talk
==========

The code presented at Sheffield Go, 7th March 2019.

## [Slides link](https://docs.google.com/presentation/d/e/2PACX-1vSXVZ2l-BYUeuQ6fLgCH5oGfKeXTsYB360Z0N3xe77WxGatqfUG2XOoOef4gzzQFJT14Ps7gaa-BOmx/pub)

To run, just say `make`.

To take a look at the output of the program, run `go run .`, or build it and run
it as you would any go program.

# What's going on?

The goal of this talk was to show you how you could use some existing tooling to make a binary from Go code.

Of course, this isn't possible in 20 minutes, so we take some liberties:

1) We use Go's x/tools/ssa (static single assignment) package, which enables us to turn Go's langauge semantics into a fairly simple data structure we can work with. (This is called the Intermediate Representation, or "IR").

2) We use LLVM, which is a compiler framework which also uses static single assignment for its IR.

3) The goal of this program then is to translate from one IR to the other.

To show some output, we use `printf`, which we steal from `libc`, C's runtime. We compile the resulting intermediate `clang`, the C frontend, which happens to be able to compile `.ll` (LLVM's "assembly" format).

This compiler is not at all general. It supports only the `+` operator and calling the `println` function (which is actually libc's printf function). So it can't do very much. But hopefully that simplicity is also what allows some newcomers to understand it.

# go2ll

In the near future I plan to publish `go2ll`, which is a slightly more sophisticated frontend. It will only ever be a toy, because, for example, it is unlikely to implement garbage collection and goroutines. This means it won't be good for abitrary Go programs. On the other hand, I can still think of a few interesting uses, such as for speeding up CPU intensive compute kernels. I have already been able to demonstrate 30-40% speedups in already fairly well tuned code such as that for computing SHA1 and `strconv.ParseFloat`.
