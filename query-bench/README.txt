lorem -p 1000000 -w 10 > list
go test -bench=.


BenchmarkRoaringScanTerm-8                    68          15659009 ns/op
BenchmarkInvertedScanTerm-8                  856           1350597 ns/op
BenchmarkRoaringScanOr-8                      72          15984605 ns/op
BenchmarkInvertedScanOr-8                     78          15072876 ns/op
BenchmarkRoaringScanAnd-8                   6458            169041 ns/op
BenchmarkInvertedScanAnd-8                   182           6499305 ns/op
BenchmarkRoaringScanAndNot-8                1190           1008909 ns/op
BenchmarkInvertedScanAndNot-8                418           2742980 ns/op
BenchmarkRoaringScanAndCompex-8             1052           1099928 ns/op
BenchmarkInvertedScanAndCompex-8             124           9630326 ns/op