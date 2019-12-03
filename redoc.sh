
echo "# github.com/jackdoe/go-query []int32 query interface" > README.md
echo '---' >> README.md
echo '[![Build Status](https://travis-ci.org/jackdoe/go-query.svg?branch=master)](https://travis-ci.org/jackdoe/go-query) [![codecov](https://codecov.io/gh/jackdoe/go-query/branch/master/graph/badge.svg)](https://codecov.io/gh/jackdoe/go-query)' >> README.md
echo '---' >> README.md
godocdown github.com/jackdoe/go-query/util/tokenize > util/tokenize/README.md
godocdown github.com/jackdoe/go-query/util/norm > util/norm/README.md
godocdown github.com/jackdoe/go-query/util/index > util/index/README.md
godocdown github.com/jackdoe/go-query/util/spec > util/spec/README.md
godocdown github.com/jackdoe/go-query/util > util/README.md
godocdown github.com/jackdoe/go-query >> README.md

echo '---' >> README.md

cat util/README.md >> README.md

echo '---' >> README.md

cat util/index/README.md >> README.md