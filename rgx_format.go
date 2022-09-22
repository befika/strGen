package strGen

import (
	"fmt"
	"math"
	"os"
	"regexp/syntax"
)

// opToString gets the string name of a regular expression operation.
func (g *Generator) OpToString(s *state, re *syntax.Regexp) string {
	op := re.Op
	switch op {
	case syntax.OpNoMatch:
	case syntax.OpEmptyMatch:
		return ""
	case syntax.OpLiteral:
		res := ""
		for _, r := range re.Rune {
			res += string(r)
		}
		return res
	case syntax.OpCharClass:
		// number of possible chars
		sum := 0
		for i := 0; i < len(re.Rune); i += 2 {
			if g.debug {
				fmt.Printf("Range: %#U-%#U\n", re.Rune[i], re.Rune[i+1])
			}
			sum += int(re.Rune[i+1]-re.Rune[i]) + 1
			if re.Rune[i+1] == runeRangeEnd {
				sum = -1
				break
			}
		}
		// pick random char in range (inverse match group)
		if sum == -1 {
			possibleChars := []uint8{}
			for j := 0; j < len(printableChars); j++ {
				c := printableChars[j]
				//fmt.Printf("Char %c %d\n", c, c)
				// Check c in range
				for i := 0; i < len(re.Rune); i += 2 {
					if rune(c) >= re.Rune[i] && rune(c) <= re.Rune[i+1] {
						possibleChars = append(possibleChars, c)
						break
					}
				}
			}
			//fmt.Println("Possible chars: ", possibleChars)
			if len(possibleChars) > 0 {
				c := possibleChars[g.rand.Intn(len(possibleChars))]
				if g.debug {
					fmt.Printf("Generated rune %c for inverse range %v\n", c, re)
				}
				return string([]byte{c})
			}
		}
		if g.debug {
			fmt.Println("Char range: ", sum)
		}
		r := g.rand.Intn(int(sum))
		var ru rune
		sum = 0
		for i := 0; i < len(re.Rune); i += 2 {
			gap := int(re.Rune[i+1]-re.Rune[i]) + 1
			if sum+gap > r {
				ru = re.Rune[i] + rune(r-sum)
				break
			}
			sum += gap
		}
		if g.debug {
			fmt.Printf("Generated rune %c for range %v\n", ru, re)
		}
		return string(ru)
	case syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		chars := printableChars
		if op == syntax.OpAnyCharNotNL {
			chars = printableCharsNoNL
		}
		c := chars[g.rand.Intn(len(chars))]
		return string([]byte{c})
	case syntax.OpBeginLine:
	case syntax.OpEndLine:
	case syntax.OpBeginText:
	case syntax.OpEndText:
	case syntax.OpWordBoundary:
	case syntax.OpNoWordBoundary:
	case syntax.OpCapture:
		if g.debug {
			fmt.Println("OpCapture", re.Sub, len(re.Sub))
		}
		return g.OpToString(s, re.Sub0[0])
	case syntax.OpStar:
		// Repeat zero or more times
		res := ""
		count := g.rand.Intn(s.limit + 1)
		for i := 0; i < count; i++ {
			for _, r := range re.Sub {
				res += g.OpToString(s, r)
			}
		}
		return res
	case syntax.OpPlus:
		// Repeat one or more times
		res := ""
		count := g.rand.Intn(s.limit) + 1
		for i := 0; i < count; i++ {
			for _, r := range re.Sub {
				res += g.OpToString(s, r)
			}
		}
		return res
	case syntax.OpQuest:
		// Zero or one instances
		res := ""
		count := g.rand.Intn(2)
		if g.debug {
			fmt.Println("Quest", count)
		}
		for i := 0; i < count; i++ {
			for _, r := range re.Sub {
				res += g.OpToString(s, r)
			}
		}
		return res
	case syntax.OpRepeat:
		// Repeat one or more times
		if g.debug {
			fmt.Println("OpRepeat", re.Min, re.Max)
		}
		res := ""
		count := 0
		re.Max = int(math.Min(float64(re.Max), float64(s.limit)))
		if re.Max > re.Min {
			count = g.rand.Intn(re.Max - re.Min + 1)
		}
		if g.debug {
			fmt.Println(re.Max, count)
		}
		for i := 0; i < re.Min || i < (re.Min+count); i++ {
			for _, r := range re.Sub {
				res += g.OpToString(s, r)
			}
		}
		return res
	case syntax.OpConcat:
		// Concatenate sub-regexes
		res := ""
		for _, r := range re.Sub {
			res += g.OpToString(s, r)
		}
		return res
	case syntax.OpAlternate:
		if g.debug {
			fmt.Println("OpAlternative", re.Sub, len(re.Sub))
		}
		i := g.rand.Intn(len(re.Sub))
		return g.OpToString(s, re.Sub[i])
	default:
		fmt.Fprintln(os.Stderr, "[reg-gen] Unhandled op: ", op)
	}
	return ""
}
