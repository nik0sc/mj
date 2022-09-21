# mj

## What is this?

mj is a Mahjong-solving library. Currently it can tell you the best grouping of tiles in a hand as well as some winning conditions.

You can try [handcheck](https://nik0sc.github.io/handcheck/) in your browser. Source code for the Go WASM entry point is at [cmd/handcheck_wasm_tinygo](https://github.com/nik0sc/mj/tree/master/cmd/handcheck_wasm_tinygo) and the Web frontend at [assets_tinygo](https://github.com/nik0sc/mj/tree/master/assets_tinygo).

It targets Go 1.16 and builds under official Go and tinygo. Tinygo is recommended for the wasm version due to smaller code size.

## Benchmarking

Performance of the checkers on certain degenerate hands.

- HandRLE: the run-length encoded hand representation
- Count: the map[Tile]int hand representation
- AllP: `b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1 b1` (cannot occur in the course of the game)
- AllPReal: `b1 b1 b1 b2 b2 b2 b3 b3 b3 b4 b4 b4 b5 b5` (can occur)
- AllC: `b1 b2 b3 b3 b4 b5 b5 b6 b7 b7 b8 b9 b9 b9`
- NS: `w1 b7 w4 c5 b9 he w5 hf w5 c3 b8 hf hn hf` (plausible hand)

HandRLE is faster. Maps are not the best way to represent the small number of tiles in each hand.

```
HandRLE: range Entries

Benchmark_OptHandRLEChecker_AllP-12                96655             11118 ns/op            4669 B/op        217 allocs/op
Benchmark_OptHandRLEChecker_AllPReal-12             2239            561164 ns/op          188182 B/op       8412 allocs/op
Benchmark_OptHandRLEChecker_AllC-12                 4593            262267 ns/op           93175 B/op       3743 allocs/op
Benchmark_OptHandRLEChecker_NS-12                  66718             17722 ns/op            6702 B/op        232 allocs/op
PASS
ok      github.com/nik0sc/mj/handcheck  5.618s

HandRLE: ForEach

Benchmark_OptHandRLEChecker_AllP-12                99618             11184 ns/op            4620 B/op        204 allocs/op
Benchmark_OptHandRLEChecker_AllPReal-12             2235            534924 ns/op          183853 B/op       8115 allocs/op
Benchmark_OptHandRLEChecker_AllC-12                 4620            249328 ns/op           89256 B/op       3580 allocs/op
Benchmark_OptHandRLEChecker_NS-12                  70348             16978 ns/op            6230 B/op        220 allocs/op
PASS
ok      github.com/nik0sc/mj/handcheck  5.433s

Counter: range Entries

Benchmark_OptCountChecker_AllP-12                  52147             21158 ns/op            9661 B/op        328 allocs/op
Benchmark_OptCountChecker_AllPReal-12                864           1382612 ns/op          417036 B/op      14509 allocs/op
Benchmark_OptCountChecker_AllC-12                   1599            755546 ns/op          190976 B/op       6313 allocs/op
Benchmark_OptCountChecker_NS-12                    21586             52459 ns/op           11393 B/op        344 allocs/op
PASS
ok      github.com/nik0sc/mj/handcheck  6.039s

Counter: ForEach

Benchmark_OptCountChecker_AllP-12                  53492             21314 ns/op            9660 B/op        328 allocs/op
Benchmark_OptCountChecker_AllPReal-12                874           1387995 ns/op          417039 B/op      14509 allocs/op
Benchmark_OptCountChecker_AllC-12                   1570            761786 ns/op          190933 B/op       6312 allocs/op
Benchmark_OptCountChecker_NS-12                    22988             51988 ns/op           11393 B/op        344 allocs/op
PASS
ok      github.com/nik0sc/mj/handcheck  5.797s
```
