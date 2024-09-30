package configkit

import (
	"context"
	"io"
	"sort"
	"strings"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/enginekit/collection/sets"
	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
)

// ToString returns a human-readable string representation of the given entity.
func ToString(e Entity) string {
	var b strings.Builder

	if err := e.AcceptVisitor(
		context.Background(),
		&stringer{w: &b},
	); err != nil {
		panic(err)
	}

	return b.String()
}

// stringer is a Visitor that builds string representations of entities.
type stringer struct {
	w io.Writer
}

func (s *stringer) VisitApplication(ctx context.Context, cfg Application) error {
	id := cfg.Identity()

	must.Fprintf(
		s.w,
		"application %s (%s) %s\n",
		id.Name,
		id.Key,
		cfg.TypeName(),
	)

	v := &stringer{
		w: indent.NewIndenter(s.w, nil),
	}

	for _, h := range sortHandlers(cfg.Handlers()) {
		must.WriteByte(s.w, '\n')
		must.WriteString(v.w, "- ")

		if err := h.AcceptVisitor(ctx, v); err != nil {
			return err
		}
	}

	return nil
}

func (s *stringer) visitHandler(cfg Handler) error {
	id := cfg.Identity()

	var flags []string
	if cfg.IsDisabled() {
		flags = append(flags, "disabled")
	}

	flagString := ""
	if len(flags) > 0 {
		flagString = " [" + strings.Join(flags, ", ") + "]"
	}

	must.Fprintf(
		s.w,
		"%s %s (%s) %s%s\n",
		cfg.HandlerType(),
		id.Name,
		id.Key,
		cfg.TypeName(),
		flagString,
	)

	names := cfg.MessageNames()

	for _, p := range sortNameKinds(names.Kinds, names.Consumed) {
		if p.Kind == message.TimeoutKind {
			break
		}

		must.Fprintf(
			s.w,
			"    handles %s%s\n",
			p.Name,
			p.Kind.Symbol(),
		)
	}

	for _, p := range sortNameKinds(names.Kinds, names.Produced) {
		must.Fprintf(
			s.w,
			"    %s %s%s\n",
			message.MapKind(p.Kind, "executes", "records", "schedules"),
			p.Name,
			p.Kind.Symbol(),
		)
	}

	return nil
}

func (s *stringer) VisitAggregate(_ context.Context, cfg Aggregate) error {
	return s.visitHandler(cfg)
}

func (s *stringer) VisitProcess(_ context.Context, cfg Process) error {
	return s.visitHandler(cfg)
}

func (s *stringer) VisitIntegration(_ context.Context, cfg Integration) error {
	return s.visitHandler(cfg)
}

func (s *stringer) VisitProjection(_ context.Context, cfg Projection) error {
	return s.visitHandler(cfg)
}

// sortHandlers returns a set of handlers sorted by their name.
func sortHandlers(handlers HandlerSet) []Handler {
	sorted := make([]Handler, 0, len(handlers))

	for _, h := range handlers {
		sorted = append(sorted, h)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Identity().Name < sorted[j].Identity().Name
	})

	return sorted
}

type pair struct {
	Name message.Name
	Kind message.Kind
}

// sortNameKinds returns a list of name/kind pairs, sorted by name.
// Timeout messages are always sorted towards the end.
func sortNameKinds(
	kinds map[message.Name]message.Kind,
	names sets.Set[message.Name],
) []pair {
	sorted := make([]pair, 0, names.Len())

	for n := range names.All() {
		sorted = append(sorted, pair{n, kinds[n]})
	}

	sort.Slice(sorted, func(i, j int) bool {
		pi := sorted[i]
		pj := sorted[j]

		if pi.Kind == message.TimeoutKind && pj.Kind != message.TimeoutKind {
			return false
		}

		return pi.Name.String() < pj.Name.String()
	})

	return sorted
}
