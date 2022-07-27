package dml

import (
	"os"
	"regexp"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/dogmatiq/configkit/dml/ast"
)

func ParseFile(filename string) (*ast.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parser.Parse(filename, f)
}

var parser = participle.MustBuild[ast.File](
	participle.Elide("Comment", "Whitespace"),
	participle.Lexer(
		lexer.MustSimple(
			[]lexer.SimpleRule{
				{
					Name:    `Ident`,
					Pattern: `[A-Z][a-zA-Z0-9]*`,
				},
				{
					Name:    `UUID`,
					Pattern: `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`,
				},
				{
					Name:    `Keyword`,
					Pattern: `[a-z]+`,
				},
				{
					Name:    `Operator`,
					Pattern: `[` + regexp.QuoteMeta(`?![].()`) + `]|->`,
				},
				{
					Name:    `Comment`,
					Pattern: `//[^\n]*`,
				},
				{
					Name:    `Documentation`,
					Pattern: `#[^\n]*\n`,
				},
				{
					Name:    `Whitespace`,
					Pattern: `[ \t\n\r]+`,
				},
			},
		),
	),
)
