lorem -p 1000000 -w 10 > list
go test -bench=.

# 1m docs
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


# 1m docs
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


# 100k docs

BenchmarkBleveScanAndTwo-8                    14          90464792 ns/op
BenchmarkRoaringScanAndTwo-8               11078            104069 ns/op
BenchmarkInvertedScanAndTwo-8               3675            325290 ns/op
BenchmarkRoaringScanAndOne-8                1156           1137950 ns/op
BenchmarkInvertedScanAndOne-8               1544            753181 ns/op
BenchmarkRoaringScanTerm-8                  1081           1157883 ns/op
BenchmarkInvertedScanTerm-8                 8043            140601 ns/op
BenchmarkRoaringScanOr-8                    1075           1066297 ns/op
BenchmarkInvertedScanOr-8                    730           1577617 ns/op
BenchmarkRoaringScanAnd-8                  31735             43779 ns/op
BenchmarkInvertedScanAnd-8                  3415            343475 ns/op
BenchmarkRoaringScanAndNot-8                9236            125304 ns/op
BenchmarkInvertedScanAndNot-8               3501            341491 ns/op
BenchmarkRoaringScanAndCompex-8             6504            157004 ns/op
BenchmarkInvertedScanAndCompex-8            1149            932140 ns/op