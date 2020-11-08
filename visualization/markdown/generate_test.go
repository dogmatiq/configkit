package markdown_test

import (
	"strings"
	"testing"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/visualization/markdown"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

func TestGenerate_coverage(t *testing.T) {
	app := &fixtures.Application{
		ConfigureFunc: func(c dogma.ApplicationConfigurer) {
			c.Identity("app", "a07d0caf-d9d0-4f9f-97d3-8779bcc304ab")

			c.RegisterAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("aggregate", "9032ca75-d734-49c6-9466-1fe17a9db1fa")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageE{})
					c.ProducesEventType(fixtures.MessageF{})
				},
			})

			c.RegisterProcess(&fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("process", "66ee8f6d-2a03-42a1-a24a-ce82f370b556")
					c.ConsumesEventType(fixtures.MessageE{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.ProducesCommandType(fixtures.MessageX{}) // not handled by this app
				},
			})

			c.RegisterIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("integration", "09ca669d-e27a-4221-ab84-d7d0f241c902")
					c.ConsumesCommandType(fixtures.MessageI{})
				},
			})

			c.RegisterProjection(&fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("projection", "c2d36b03-dab4-4761-8bd6-efc47077c458")
					c.ConsumesEventType(fixtures.MessageE{})
					c.ConsumesEventType(fixtures.MessageF{})
					c.ConsumesEventType(fixtures.MessageY{}) // not produced by this app
				},
			})
		},
	}

	w := &strings.Builder{}
	cfg := configkit.FromApplication(app)

	if err := Generate(w, cfg); err != nil {
		t.Fatal(err)
	}
}
