# rewindomata
New side of automaton

## Installation
The package depends on `github.com/golang-collections/collections` and `github.com/uber-go/atomic`, so install it first, or if you have `dep` installed, `ensure` the repository after cloning or downloading.

For testing run: `go test -run=XXX -bench=. -benchmem > testFile.txt`

## The problem

Need to write a package that will parse regular expressions to non-deterministic finite state machines (see [the repo](https://github.com/PaulRaUnite/dsl_lab1)) and a package with the following types of word checking algorithms:

- pair based: queue of pairs of form `(state, word part)` used by algorithm that for a given pair try to get next states using first letter of word, if there is some next states, it adds pairs `(new state, word part without first letter)` to the queue; if final state found the algorithm finishes its work, otherwise works while the queue would be empty;
- the above algorithm, but processing is parallel;
- the above algorithm, but "pop"ing from the queue is random.

The idea of the problem is to try to understand, can the following variants by faster than original.

## Implementation

`ast` package provides AST(abstract syntax tree) type and parser for suffix notation variant of regular expression(for simplicity and thus faster development of parser).

`automaton` package provides `Acceptor` type, that is a non-deterministic finite state machine, utility components as queue, AST transforming constructor and word recognition algorithms.

## About parallel algorithm

![Parallel](./parallel.svg)

Data counter counts amount of payload in queue and if it is 0, it means that there is nothing to do.
Algorithm starts pool of workers, that try to recognise a few words by the automaton, and returns resulting array.

## Results

| Benchmark name | loop steps | speed: less better | memory consumption: less better | allocations |
|-------|-------|-------|-------|-------|
| BenchmarkAcceptor_FrontSearchPositive       |     	   10000	|    152007 ns/op	|   11033 B/op	|     240 allocs/op |
| BenchmarkAcceptor_AtomicSearchPositive          | 	   10000	|    100185 ns/op	|   19037 B/op	 |     43 allocs/op |
| BenchmarkAcceptor_AtomicParallelSearchPositive   |	    2000	|    733733 ns/op	|  48591 B/op	 |     84 allocs/op |
| BenchmarkAcceptor_StochasticSearchPositive      |	   10000	 |   126136 ns/op	 |  17704 B/op	 |     42 allocs/op |
| BenchmarkAcceptor_FrontSearchNegative      |      	   10000	|    126630 ns/op	|    8665 B/op	 |    182 allocs/op |
| BenchmarkAcceptor_AtomicSearchNegative |           	   20000	 |    83764 ns/op	|   16091 B/op	|  
    35 allocs/op |
| BenchmarkAcceptor_AtomicParallelSearchNegative   	|    2000	|    642846 ns/op	|   45381 B/op	  |    84 allocs/op |
| BenchmarkAcceptor_StochasticSearchNegative   |    	   10000	|    114276 ns/op	|   14736 B/op | 35 allocs/op |

Maybe, my test cases are not representative, but I guess in case of parallel search nothing much can be changed, it is not effective at all, because of small amount of useful computations in workers: a lot of time they spend synchronizing with each other. 