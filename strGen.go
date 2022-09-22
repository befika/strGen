// Package reggen generates text based on regex definitions
package strGen

import (
	"math/rand"
	"regexp/syntax"
	"time"
)

const runeRangeEnd = 0x10ffff
const printableChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ \t\n\r"

var printableCharsNoNL = printableChars[:len(printableChars)-2]

type state struct {
	limit int
}

type Generator struct {
	re    *syntax.Regexp
	rand  *rand.Rand
	debug bool
}

// limit is the maximum number of times star, range or plus should repeat
// i.e. [0-9]+ will generate at most 10 characters if this is set to 10
func (g *Generator) Generate(limit int) string {
	return g.OpToString(&state{limit: limit}, g.re)
}

// create a new generator
func NewGenerator(regex string) (*Generator, error) {

	// parse the rgx pattern using regexp/syntax package
	re, err := syntax.Parse(regex, syntax.Perl)
	if err != nil {
		return nil, err
	}
	//fmt.Println("Compiled re ", re)
	return &Generator{
		re:   re,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (gen *Generator) SetSeed(seed int64) {
	gen.rand = rand.New(rand.NewSource(seed))
}

func Generate(regex string, limit int) (string, error) {
	g, err := NewGenerator(regex)
	if err != nil {
		return "", err
	}
	return g.Generate(limit), nil
}
