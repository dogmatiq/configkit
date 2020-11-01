package dot_test

import (
	"testing"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/visualization/dot"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
)

func TestGenerate_coverage(t *testing.T) {
	app := &fixtures.Application{
		ConfigureFunc: func(c dogma.ApplicationConfigurer) {
			c.Identity("app", "a07d0caf-d9d0-4f9f-97d3-8779bcc304ab")

			c.RegisterAggregate(&fixtures.AggregateMessageHandler{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("aggregate", "b2a8b880-5a1a-4792-ab03-5675b002230a")
					c.ConsumesCommandType(fixtures.MessageC{})
					c.ProducesEventType(fixtures.MessageE{})
					c.ProducesEventType(fixtures.MessageF{})
				},
			})

			c.RegisterProcess(&fixtures.ProcessMessageHandler{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("process", "3d5bb944-1cb7-40f4-9298-e154acd5effd")
					c.ConsumesEventType(fixtures.MessageE{})
					c.ProducesCommandType(fixtures.MessageC{})
					c.ProducesCommandType(fixtures.MessageX{}) // not handled by this app
				},
			})

			c.RegisterIntegration(&fixtures.IntegrationMessageHandler{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("integration", "5a496ba8-92f4-439e-bdba-d0e4ef6dd03d")
					c.ConsumesCommandType(fixtures.MessageI{})
				},
			})

			c.RegisterProjection(&fixtures.ProjectionMessageHandler{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("projection", "3f060ff7-630a-4446-8313-35ace689d5ce")
					c.ConsumesEventType(fixtures.MessageE{})
					c.ConsumesEventType(fixtures.MessageF{})
					c.ConsumesEventType(fixtures.MessageY{}) // not produced by this app
				},
			})
		},
	}

	cfg := configkit.FromApplication(app)

	_, err := Generate(cfg)
	if err != nil {
		t.Fatal(err)
	}
}
