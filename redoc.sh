
cp README.header README.md

godocdown github.com/rekki/go-query/util/tokenize > util/tokenize/README.md
godocdown github.com/rekki/go-query/util/norm > util/norm/README.md
godocdown github.com/rekki/go-query/util/index > util/index/README.md
godocdown github.com/rekki/go-query/util/go_query_dsl > util/go_query_dsl/README.md
godocdown github.com/rekki/go-query/util > util/README.md
godocdown github.com/rekki/go-query >> README.md

echo '---' >> README.md

cat util/README.md >> README.md

echo '---' >> README.md

cat util/index/README.md >> README.md