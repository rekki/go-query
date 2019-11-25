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


--- after Mon 25 Nov 23:27:51 CET 2019 ---


BenchmarkRoaringScanAndTwo-8                1176           1037217 ns/op
BenchmarkInvertedScanAndTwo-8                387           3090059 ns/op
BenchmarkRoaringScanAndOne-8                  78          15943453 ns/op
BenchmarkInvertedScanAndOne-8                176           6793943 ns/op
BenchmarkRoaringScanTerm-8                    76          15927295 ns/op
BenchmarkInvertedScanTerm-8                  896           1326020 ns/op
BenchmarkRoaringScanOr-8                      70          15767382 ns/op
BenchmarkInvertedScanOr-8                     80          14824888 ns/op
BenchmarkRoaringScanAnd-8                   6708            171946 ns/op
BenchmarkInvertedScanAnd-8                   372           3244076 ns/op
BenchmarkRoaringScanAndNot-8                1125           1017742 ns/op
BenchmarkInvertedScanAndNot-8                363           3246193 ns/op
BenchmarkRoaringScanAndCompex-8             1059           1100383 ns/op
BenchmarkInvertedScanAndCompex-8             142           8334875 ns/op
