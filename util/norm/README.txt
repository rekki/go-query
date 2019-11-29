package norm // import "github.com/jackdoe/go-query/util/norm"

Simlpe normalizer chain

Example:

    package main
    import n "github.com/jackdoe/go-query/util/norm"
    func main() {
    	nor := []n.Normalizer{n.NewUnaccent(), n.NewLowerCase(), n.NewSpaceBetweenDigits(), n.NewCleanup(n.BASIC_NON_ALPHANUMERIC),n.NewTrim(" ")}
    	normal := n.Normalize("Hęllö wÖrld. べぺ Ł2ł  ", nor...)

    	fmt.Printf("%s",normal) // hello world へへ l 2 l
    }

VARIABLES

var BASIC_NON_ALPHANUMERIC = regexp.MustCompile(`[^\pL\pN]+`)

FUNCTIONS

func Normalize(s string, normalizers ...Normalizer) string

TYPES

type Cleanup struct {
	// Has unexported fields.
}

func NewCleanup(re *regexp.Regexp) *Cleanup

func (l *Cleanup) Apply(s string) string

type Custom struct {
	// Has unexported fields.
}

func NewCustom(f func(string) string) *Custom

func (l *Custom) Apply(s string) string

type LowerCase struct{}

func NewLowerCase() *LowerCase

func (l *LowerCase) Apply(s string) string

type Normalizer interface {
	Apply(string) string
}

type SpaceBetweenDigits struct{}

func NewSpaceBetweenDigits() *SpaceBetweenDigits

func (l *SpaceBetweenDigits) Apply(s string) string

type Trim struct {
	// Has unexported fields.
}

func NewTrim(cutset string) *Trim

func (l *Trim) Apply(s string) string

type Unaccent struct{}

func NewUnaccent() *Unaccent

func (l *Unaccent) Apply(s string) string

