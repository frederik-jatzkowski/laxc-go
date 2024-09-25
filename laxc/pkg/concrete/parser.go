package concrete

import (
	"fmt"
	"laxc/pkg/lex"

	"github.com/alecthomas/participle/v2"
)

var parser = participle.MustBuild[Prog](
	participle.Lexer(lex.Lexer),
	participle.Elide("Whitespace", "Comment"),
	participle.UseLookahead(1),
)

func Parse(fileName, source string) (Prog, error) {
	prog, err := parser.ParseString(fileName, source)
	if err != nil {
		return Prog{}, err
	}

	if prog == nil {
		return Prog{}, fmt.Errorf("parser did not return an abstract")
	}

	return *prog, err
}
