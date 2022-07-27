package ast

type File struct {
	App      *App       `parser:"@@"`
	Handlers []*Handler `parser:"@@*"`
}

type App struct {
	Name string  `parser:"'app' @Ident"`
	Key  *string `parser:"( '(' @UUID ')' )?"`
}

type Handler struct {
	Aggregate   *Aggregate   `parser:"  @@"`
	Process     *Process     `parser:"| @@"`
	Integration *Integration `parser:"| @@"`
	Projection  *Projection  `parser:"| @@"`
}

type Aggregate struct {
	Name string  `parser:"'aggregate' @Ident"`
	Key  *string `parser:"( '(' @UUID ')' )?"`
}

type Process struct {
	Name string  `parser:"'process' @Ident"`
	Key  *string `parser:"( '(' @UUID ')' )?"`
}

type Integration struct {
	Name string  `parser:"'integration' @Ident"`
	Key  *string `parser:"( '(' @UUID ')' )?"`
}

type Projection struct {
	Name string  `parser:"'projection' @Ident"`
	Key  *string `parser:"( '(' @UUID ')' )?"`
}
