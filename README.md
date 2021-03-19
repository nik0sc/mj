# mj

## Benchmarking

### HandRLE: range Entries
```
Benchmark_OptHandRLEChecker_AllP-12                96655             11118 ns/op            4669 B/op        217 allocs/op
Benchmark_OptHandRLEChecker_AllPReal-12             2239            561164 ns/op          188182 B/op       8412 allocs/op
Benchmark_OptHandRLEChecker_AllC-12                 4593            262267 ns/op           93175 B/op       3743 allocs/op
Benchmark_OptHandRLEChecker_NS-12                  66718             17722 ns/op            6702 B/op        232 allocs/op
PASS
ok      github.com/nik0sc/mj/handcheck  5.618s
```

### HandRLE: ForEach
```
Benchmark_OptHandRLEChecker_AllP-12                99618             11184 ns/op            4620 B/op        204 allocs/op
Benchmark_OptHandRLEChecker_AllPReal-12             2235            534924 ns/op          183853 B/op       8115 allocs/op
Benchmark_OptHandRLEChecker_AllC-12                 4620            249328 ns/op           89256 B/op       3580 allocs/op
Benchmark_OptHandRLEChecker_NS-12                  70348             16978 ns/op            6230 B/op        220 allocs/op
PASS
ok      github.com/nik0sc/mj/handcheck  5.433s
```

### Counter: range Entries
```
Benchmark_OptCountChecker_AllP-12                  52147             21158 ns/op            9661 B/op        328 allocs/op
Benchmark_OptCountChecker_AllPReal-12                864           1382612 ns/op          417036 B/op      14509 allocs/op
Benchmark_OptCountChecker_AllC-12                   1599            755546 ns/op          190976 B/op       6313 allocs/op
Benchmark_OptCountChecker_NS-12                    21586             52459 ns/op           11393 B/op        344 allocs/op
PASS
ok      github.com/nik0sc/mj/handcheck  6.039s
```
### Counter: ForEach
```
Benchmark_OptCountChecker_AllP-12                  53492             21314 ns/op            9660 B/op        328 allocs/op
Benchmark_OptCountChecker_AllPReal-12                874           1387995 ns/op          417039 B/op      14509 allocs/op
Benchmark_OptCountChecker_AllC-12                   1570            761786 ns/op          190933 B/op       6312 allocs/op
Benchmark_OptCountChecker_NS-12                    22988             51988 ns/op           11393 B/op        344 allocs/op
PASS
ok      github.com/nik0sc/mj/handcheck  5.797s
```