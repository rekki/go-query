# util
--
Simlpe utils to tokenize and normalize text

Example: package main

import (

    "fmt"

    n "github.com/jackdoe/go-query/util/norm"
    t "github.com/jackdoe/go-query/util/tokenize"

)

func main() {

    tokenizer := []t.Tokenizer{
    	t.NewWhitespace(),
    	t.NewLeftEdge(1),
    	t.NewUnique(),
    }
    normalizer := []n.Normalizer{
    	n.NewUnaccent(),
    	n.NewLowerCase(),
    	n.NewSpaceBetweenDigits(),
    	n.NewCleanup(n.BASIC_NON_ALPHANUMERIC),
    	n.NewTrim(" "),
    }

    tokens := t.Tokenize(
    	n.Normalize("Hęllö World yęar2019 ", normalizer...),
    	tokenizer...,
    )

    fmt.Printf("%v", tokens)
    // prints [h he hel hell hello w wo wor worl world y ye yea year 2 20 201 2019]

}
