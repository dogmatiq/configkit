package configkit

import (
	"context"
	"io"
	"sort"
	"strings"

	"github.com/dogmatiq/configkit/message"
	"github.com/dogmatiq/iago/indent"
	"github.com/dogmatiq/iago/must"
)

// ToString returns a human-readable string representation of the given entity.
func ToString(e Entity) string {
	var b strings.Builder
	e.AcceptVisitor(nil, &stringer{w: &b})
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
		h.AcceptVisitor(ctx, v)
	}

	return nil
}

func (s *stringer) visitHandler(cfg Handler) error {
	id := cfg.Identity()

	must.Fprintf(
		s.w,
		"%s %s (%s) %s\n",
		cfg.HandlerType(),
		id.Name,
		id.Key,
		cfg.TypeName(),
	)

	for _, p := range sortNameRoles(cfg.MessageNames().Consumed) {
		if p.Role == message.TimeoutRole {
			break
		}

		must.Fprintf(
			s.w,
			"    consumes %s%s\n",
			p.Name,
			p.Role.Marker(),
		)
	}

	for _, p := range sortNameRoles(cfg.MessageNames().Produced) {
		verb := "produces"
		if p.Role == message.TimeoutRole {
			verb = "schedules"
		}

		must.Fprintf(
			s.w,
			"    %s %s%s\n",
			verb,
			p.Name,
			p.Role.Marker(),
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
	Role message.Role
}

// sortNameRoles returns a list of name/role pairs, sorted by name.
// Timeout messages are always sorted towards the end.
func sortNameRoles(names message.NameRoles) []pair {
	sorted := make([]pair, 0, len(names))

	for n, r := range names {
		sorted = append(sorted, pair{n, r})
	}

	sort.Slice(sorted, func(i, j int) bool {
		pi := sorted[i]
		pj := sorted[j]

		if pi.Role == message.TimeoutRole && pj.Role != message.TimeoutRole {
			return false
		}

		return pi.Name.String() < pj.Name.String()
	})

	return sorted
}
