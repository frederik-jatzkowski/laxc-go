package lex

import "github.com/alecthomas/participle/v2/lexer"

// var simple = lexer.MustSimple([]lexer.SimpleRule{
// 	{Name: "IntLit", Pattern: `[0-9]+`},
// 	{Name: "Keyword", Pattern: `(declare|begin|is|end|div|mod)`},
// 	{Name: "Special", Pattern: `(\+|-|\*|:|;|\(|\))`},
// 	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
// 	{Name: "Comment", Pattern: `##[^(\*\))]*##`},
// 	{Name: "Ident", Pattern: `[a-zA-Z]+`},
// })

var stateful = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "RealLit", Pattern: `([0-9]+e[\+|-]?[0-9]+)|([0-9]+\.[0-9]+(e[\+|-]?[0-9]+)?)`},
		{Name: "IntLit", Pattern: `[0-9]+`},
		{Name: "Keyword", Pattern: `(declare|begin|is|end|if|then|else|not|//|case|of)`},
		{Name: "Assign", Pattern: `:=`},
		{Name: "AddOp", Pattern: `(\+|-)`},
		{Name: "MulOp", Pattern: `(\*|div|mod|/)`},
		{Name: "Comment", Pattern: `\(\*([^\*]*(\*[^\)\*])?)*\*+\)`},
		{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
		{Name: "Special", Pattern: `(:|;|\(|\)|<|>|=)`},
		{Name: "Ident", Pattern: `[a-z][_a-z0-9]*`},
	},
})

var Lexer = stateful
