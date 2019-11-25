lorem -p 1000000 -w 10 > list
go test -bench=.



goos: linux
goarch: amd64
BenchmarkRoaringScanTerm-8                    72          16001046 ns/op
BenchmarkInvertedScanTerm-8                  880           1333437 ns/op
BenchmarkRoaringScanOr-8                      74          15585116 ns/op
BenchmarkInvertedScanOr-8                     73          14906093 ns/op
BenchmarkRoaringScanAnd-8                   7086            168726 ns/op
BenchmarkInvertedScanAnd-8                   186           6317747 ns/op
BenchmarkRoaringScanAndNot-8                1204           1004954 ns/op
BenchmarkInvertedScanAndNot-8                188           6376152 ns/op
BenchmarkRoaringScanAndCompex-8              964           1108790 ns/op
BenchmarkInvertedScanAndCompex-8              79          13957232 ns/op
