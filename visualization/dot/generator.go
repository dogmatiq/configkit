package dot

import (
	"sort"
	"strconv"
	"strings"

	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/enginekit/message"
	"github.com/emicklei/dot"
)

// generator generates message-flow graphs for one or more Dogma applications.
type generator struct {
	id        int
	root      *dot.Graph
	app       *dot.Graph
	kinds     map[message.Name]message.Kind
	producers map[message.Name][]dot.Node
	consumers map[message.Name][]dot.Node
}

// nextID returns a unique ID to use for a sub-graph or node.
func (g *generator) nextID() string {
	g.id++
	return strconv.Itoa(g.id)
}

// generate builds and returns a graph of the given applications.
func (g *generator) generate(apps []configkit.Application) (_ *dot.Graph, err error) {
	g.root = dot.NewGraph(dot.Directed)
	g.root.Attr("rankdir", "TB")
	g.root.Attr("concentrate", "false")
	g.root.Attr("splines", "true")
	g.root.Attr("overlap", "false")
	g.root.Attr("outputorder", "edgesfirst")

	g.kinds = map[message.Name]message.Kind{}
	g.producers = map[message.Name][]dot.Node{}
	g.consumers = map[message.Name][]dot.Node{}

	for _, app := range apps {
		g.addApp(app)
	}

	g.addForeign()

	return g.root, nil
}

// addApp adds an application to the graph as a sub-graph.
func (g *generator) addApp(cfg configkit.Application) {
	g.app = g.root.Subgraph(
		g.nextID(),
		dot.ClusterOption{},
	)
	g.app.Attr("label", "("+cfg.Identity().Name+")")
	styleApp(g.app, cfg)

	for n, m := range cfg.MessageNames() {
		g.kinds[n] = m.Kind
	}

	for _, h := range sortHandlers(cfg.Handlers()) {
		g.addHandler(h)
	}
}

// addHandler adds a handler to the graph as a node.
func (g *generator) addHandler(cfg configkit.Handler) {
	n := g.app.Node(g.nextID())
	n.Attr("label", cfg.Identity().Name)
	styleHandler(n, cfg)

	g.addEdges(cfg, n)
}

// addEdges adds edges describing the messages that are produced and consumed by
// a specific handler.
func (g *generator) addEdges(cfg configkit.Handler, n dot.Node) {
	for name, em := range cfg.MessageNames() {
		if em.IsConsumed {
			g.consumers[name] = append(g.consumers[name], n)

			for _, p := range g.producers[name] {
				g.addEdge(p, n, name)
			}
		}

		if em.IsProduced {
			g.producers[name] = append(g.producers[name], n)

			for _, c := range g.consumers[name] {
				g.addEdge(n, c, name)
			}
		}
	}
}

// addEdge adds an edge between two handler nodes.
//
// If the edge already exists, its label is expanded to include this message type.
func (g *generator) addEdge(src, dst dot.Node, n message.Name) {
	k := g.kinds[n]
	label := string(n) + k.Symbol()

	index := strings.LastIndex(label, ".")
	if index != -1 {
		label = label[index+1:]
	}

	label = " " + label + " "

	// find an existing edge and add to its label.
	for _, e := range g.root.FindEdges(src, dst) {
		labels := strings.Split(
			e.Value("label").(string),
			"\n",
		)

		labels = append(labels, label)
		sort.Strings(labels)

		e.Label(strings.Join(labels, "\n"))
		return
	}

	e := g.root.Edge(
		src,
		dst,
		label,
	)

	styleMessageEdge(e, k)
}

// addForeign adds nodes that represent producers and consumers that are
// external to any of the applications in the graph.
func (g *generator) addForeign() {
	for name, nodes := range g.producers {
		if _, ok := g.consumers[name]; ok {
			continue
		}

		for _, n := range nodes {
			g.addEdge(
				n,
				g.foreignConsumer(g.kinds[name]),
				name,
			)
		}
	}

	for name, nodes := range g.consumers {
		if _, ok := g.producers[name]; ok {
			continue
		}

		for _, n := range nodes {
			g.addEdge(
				g.foreignProducer(g.kinds[name]),
				n,
				name,
			)
		}
	}
}

// foreignConsumer adds and returns a node representing an external consumer.
func (g *generator) foreignConsumer(k message.Kind) dot.Node {
	n := g.root.Node("foreign\n" + k.String() + "\nconsumer")
	styleForeignNode(n, k)

	return n
}

// foreignProducer adds and returns a node representing an external producer.
func (g *generator) foreignProducer(k message.Kind) dot.Node {
	n := g.root.Node("foreign\n" + k.String() + "\nproducer")
	styleForeignNode(n, k)

	return n
}

// sortHandlers returns a set of handlers in an app, sorted by the number of
// message types they are aware of.
func sortHandlers(handlers configkit.HandlerSet) []configkit.Handler {
	var sorted []configkit.Handler

	for _, h := range handlers {
		sorted = append(sorted, h)
	}

	sort.Slice(
		sorted,
		func(i, j int) bool {
			in := sorted[i].MessageNames()
			jn := sorted[j].MessageNames()
			return len(in) < len(jn)
		},
	)

	return sorted
}
