go doc github.com/jackdoe/go-query/util/tokenize > util/tokenize/README.txt
go doc github.com/jackdoe/go-query/util/norm > util/norm/README.txt
go doc github.com/jackdoe/go-query/util > util/README.txt
go doc github.com/jackdoe/go-query > README.txt

echo >> README.txt
echo "----------------------------" >> README.txt
echo >> README.txt

cat util/README.txt >> README.txt