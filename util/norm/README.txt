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

var BASIC_NON_ALPHANUMERIC = regexp.MustCompile(`[^\pL\pN]+`)
func Normalize(s string, normalizers ...Normalizer) string
type Cleanup struct{ ... }
    func NewCleanup(re *regexp.Regexp) *Cleanup
type Custom struct{ ... }
    func NewCustom(f func(string) string) *Custom
type LowerCase struct{}
    func NewLowerCase() *LowerCase
type Normalizer interface{ ... }
type SpaceBetweenDigits struct{}
    func NewSpaceBetweenDigits() *SpaceBetweenDigits
type Trim struct{ ... }
    func NewTrim(cutset string) *Trim
type Unaccent struct{}
    func NewUnaccent() *Unaccent
