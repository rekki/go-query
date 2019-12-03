# norm
--
    import "github.com/jackdoe/go-query/util/norm"

Simlpe normalizer chain

Example:

    package main
    import n "github.com/jackdoe/go-query/util/norm"
    func main() {
    	nor := []n.Normalizer{n.NewUnaccent(), n.NewLowerCase(), n.NewSpaceBetweenDigits(), n.NewCleanup(n.BASIC_NON_ALPHANUMERIC),n.NewTrim(" ")}
    	normal := n.Normalize("Hęllö wÖrld. べぺ Ł2ł  ", nor...)

    	fmt.Printf("%s",normal) // hello world へへ l 2 l
    }

## Usage

```go
var BASIC_NON_ALPHANUMERIC = regexp.MustCompile(`[^\pL\pN]+`)
```

#### func  Normalize

```go
func Normalize(s string, normalizers ...Normalizer) string
```

#### type Cleanup

```go
type Cleanup struct {
}
```


#### func  NewCleanup

```go
func NewCleanup(re *regexp.Regexp) *Cleanup
```

#### func (*Cleanup) Apply

```go
func (l *Cleanup) Apply(s string) string
```

#### type Custom

```go
type Custom struct {
}
```


#### func  NewCustom

```go
func NewCustom(f func(string) string) *Custom
```

#### func (*Custom) Apply

```go
func (l *Custom) Apply(s string) string
```

#### type LowerCase

```go
type LowerCase struct{}
```


#### func  NewLowerCase

```go
func NewLowerCase() *LowerCase
```

#### func (*LowerCase) Apply

```go
func (l *LowerCase) Apply(s string) string
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
func (w *Noop) Apply(current string) string
```

#### type Normalizer

```go
type Normalizer interface {
	Apply(string) string
}
```


#### type SpaceBetweenDigits

```go
type SpaceBetweenDigits struct{}
```


#### func  NewSpaceBetweenDigits

```go
func NewSpaceBetweenDigits() *SpaceBetweenDigits
```

#### func (*SpaceBetweenDigits) Apply

```go
func (l *SpaceBetweenDigits) Apply(s string) string
```

#### type Trim

```go
type Trim struct {
}
```


#### func  NewTrim

```go
func NewTrim(cutset string) *Trim
```

#### func (*Trim) Apply

```go
func (l *Trim) Apply(s string) string
```

#### type Unaccent

```go
type Unaccent struct{}
```


#### func  NewUnaccent

```go
func NewUnaccent() *Unaccent
```

#### func (*Unaccent) Apply

```go
func (l *Unaccent) Apply(s string) string
```
