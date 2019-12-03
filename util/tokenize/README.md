# tokenize
--
    import "github.com/jackdoe/go-query/util/tokenize"

Simlpe tokenizer chain

Example:

    package main
    import t "github.com/jackdoe/go-query/util/tokenize"
    func main() {
    	tokenizer := []t.Tokenizer{t.NewWhitespace(), t.NewLeftEdge(1), t.NewUnique()}
    	tokens := t.Tokenize("hello world", tokenizer...)

    	fmt.Printf("%v",tokens) // [h he hel hell hello w wo wor worl world]
    }

## Usage

```go
const BASE_SOUNDEX = "0000"
```

#### func  EncodeSoundex

```go
func EncodeSoundex(word string) string
```

#### func  Tokenize

```go
func Tokenize(s string, tokenizers ...Tokenizer) []string
```

#### type CharNgram

```go
type CharNgram struct {
}
```


#### func  NewCharNgram

```go
func NewCharNgram(size int) *CharNgram
```

#### func (*CharNgram) Apply

```go
func (w *CharNgram) Apply(current []string) []string
```

#### type Custom

```go
type Custom struct {
}
```


#### func  NewCustom

```go
func NewCustom(f func([]string) []string) *Custom
```

#### func (*Custom) Apply

```go
func (l *Custom) Apply(s []string) []string
```

#### type LeftEdge

```go
type LeftEdge struct {
}
```


#### func  NewLeftEdge

```go
func NewLeftEdge(n int) *LeftEdge
```

#### func (*LeftEdge) Apply

```go
func (e *LeftEdge) Apply(current []string) []string
```

#### type Noop

```go
type Noop struct{}
```


#### func  NewNoop

```go
func NewNoop() *Noop
```

#### func (*Noop) Apply

```go
func (w *Noop) Apply(current []string) []string
```

#### type Soundex

```go
type Soundex struct {
}
```


#### func  NewSoundex

```go
func NewSoundex() *Soundex
```

#### func (*Soundex) Apply

```go
func (w *Soundex) Apply(current []string) []string
```

#### type Surround

```go
type Surround struct {
}
```

NewSurround("$").Apply([]string{"h","he","hel"}) -> []string{"$h","he","hel$"}

#### func  NewSurround

```go
func NewSurround(s string) *Surround
```

#### func (*Surround) Apply

```go
func (w *Surround) Apply(current []string) []string
```

#### type Tokenizer

```go
type Tokenizer interface {
	Apply([]string) []string
}
```


#### type Unique

```go
type Unique struct {
}
```


#### func  NewUnique

```go
func NewUnique() *Unique
```

#### func (*Unique) Apply

```go
func (w *Unique) Apply(current []string) []string
```

#### type Whitespace

```go
type Whitespace struct{}
```


#### func  NewWhitespace

```go
func NewWhitespace() *Whitespace
```

#### func (*Whitespace) Apply

```go
func (w *Whitespace) Apply(current []string) []string
```
