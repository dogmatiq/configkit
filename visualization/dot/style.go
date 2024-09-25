package dot

import (
	"github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	"github.com/emicklei/dot"
)

const fontName = "Helvetica"

var kindFontColors = map[message.Kind]string{
	message.CommandKind: "#ffffff",
	message.EventKind:   "#000000",
	message.TimeoutKind: "#ffffff",
}

var kindBackgroundColors = map[message.Kind]string{
	message.CommandKind: "#0066ff",
	message.EventKind:   "#ff6600",
	message.TimeoutKind: "#444444",
}

var handlerShapes = map[configkit.HandlerType]string{
	configkit.AggregateHandlerType:   "box",
	configkit.ProcessHandlerType:     "box",
	configkit.IntegrationHandlerType: "box",
	configkit.ProjectionHandlerType:  "cylinder",
}

var handlerFontColors = map[configkit.HandlerType]string{
	configkit.AggregateHandlerType:   "#000000",
	configkit.ProcessHandlerType:     "#ffffff",
	configkit.IntegrationHandlerType: "#000000",
	configkit.ProjectionHandlerType:  "#000000",
}

var handlerBackgroundColors = map[configkit.HandlerType]string{
	configkit.AggregateHandlerType:   "#fff89d",
	configkit.ProcessHandlerType:     "#c086c5",
	configkit.IntegrationHandlerType: "#ffde00",
	configkit.ProjectionHandlerType:  "#a3d350",
}

// styleApp applies style attributes to a sub-graph representing an application.
func styleApp(g *dot.Graph, _ configkit.Application) {
	g.Attr("fontname", fontName)
	g.Attr("fontsize", "18")
	g.Attr("labelloc", "b")
	g.Attr("fontcolor", "#444444")
	g.Attr("penwidth", "0")
	g.Attr("style", "filled")
	g.Attr("fillcolor", "#eeeeee")
}

// styleHandler applies style attributes to a graph node representing a Dogma
// message handler.
func styleHandler(n dot.Node, cfg configkit.Handler) {
	n.Attr("fontname", fontName)
	n.Attr("style", "filled")
	n.Attr("margin", "0.15")
	n.Attr("penwidth", "2")
	n.Attr("color", "#eeeeee")
	n.Attr("fontcolor", handlerFontColors[cfg.HandlerType()])
	n.Attr("fillcolor", handlerBackgroundColors[cfg.HandlerType()])
	n.Attr("shape", handlerShapes[cfg.HandlerType()])
}

// styleMessageEdge applies style attributes to an edge representing message
// flow between two handlers.
func styleMessageEdge(e dot.Edge, k message.Kind) {
	e.Attr("fontname", fontName)
	e.Attr("penwidth", "2")
	e.Attr("arrowsize", "0.75")
	e.Attr("color", kindBackgroundColors[k])
	e.Attr("fontcolor", kindBackgroundColors[k])
}

// styleForeignNode applies style attributes to a graph node representing a
// foreign message producer or consumer.
func styleForeignNode(n dot.Node, k message.Kind) {
	n.Attr("fontname", fontName)
	n.Attr("style", "filled")
	n.Attr("margin", "0.15")
	n.Attr("penwidth", "2")
	n.Attr("color", "#ffffff")
	n.Attr("fontcolor", kindFontColors[k])
	n.Attr("fillcolor", kindBackgroundColors[k])
	n.Attr("shape", "box")
}
