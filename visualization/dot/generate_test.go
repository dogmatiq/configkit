package dot_test

import (
	"testing"

	"github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/configkit/visualization/dot"
	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
)

func TestGenerate_coverage(t *testing.T) {
	app := &ApplicationStub{
		ConfigureFunc: func(c dogma.ApplicationConfigurer) {
			c.Identity("app", "a07d0caf-d9d0-4f9f-97d3-8779bcc304ab")

			c.RegisterAggregate(&AggregateMessageHandlerStub{
				ConfigureFunc: func(c dogma.AggregateConfigurer) {
					c.Identity("aggregate", "b2a8b880-5a1a-4792-ab03-5675b002230a")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeA]](),
						dogma.RecordsEvent[EventStub[TypeB]](),
					)
				},
			})

			c.RegisterProcess(&ProcessMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProcessConfigurer) {
					c.Identity("process", "3d5bb944-1cb7-40f4-9298-e154acd5effd")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeA]](),
						dogma.ExecutesCommand[CommandStub[TypeX]](), // not handled by this app
						dogma.SchedulesTimeout[TimeoutStub[TypeA]](),
					)
				},
			})

			c.RegisterIntegration(&IntegrationMessageHandlerStub{
				ConfigureFunc: func(c dogma.IntegrationConfigurer) {
					c.Identity("integration", "5a496ba8-92f4-439e-bdba-d0e4ef6dd03d")
					c.Routes(
						dogma.HandlesCommand[CommandStub[TypeB]](),
					)
				},
			})

			c.RegisterProjection(&ProjectionMessageHandlerStub{
				ConfigureFunc: func(c dogma.ProjectionConfigurer) {
					c.Identity("projection", "3f060ff7-630a-4446-8313-35ace689d5ce")
					c.Routes(
						dogma.HandlesEvent[EventStub[TypeA]](),
						dogma.HandlesEvent[EventStub[TypeB]](),
						dogma.HandlesEvent[EventStub[TypeX]](), // not produced by this app
					)
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
