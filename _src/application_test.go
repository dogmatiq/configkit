package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	configfixtures "github.com/dogmatiq/configkit/fixtures"
	"github.com/dogmatiq/dogma"
	"github.com/dogmatiq/dogma/fixtures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type RichApplication", func() {
	var (
		aggregate   *fixtures.AggregateMessageHandler
		process     *fixtures.ProcessMessageHandler
		integration *fixtures.IntegrationMessageHandler
		projection  *fixtures.ProjectionMessageHandler
		app         *fixtures.Application
	)

	BeforeEach(func() {
		aggregate = &fixtures.AggregateMessageHandler{
			ConfigureFunc: func(c dogma.AggregateConfigurer) {
				c.Identity("<aggregate>", "<aggregate-key>")
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			},
		}

		process = &fixtures.ProcessMessageHandler{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("<process>", "<process-key>")
				c.ConsumesEventType(fixtures.MessageB{})
				c.ConsumesEventType(fixtures.MessageE{}) // shared with <projection>
				c.ProducesCommandType(fixtures.MessageC{})
				c.SchedulesTimeoutType(fixtures.MessageT{})
			},
		}

		integration = &fixtures.IntegrationMessageHandler{
			ConfigureFunc: func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration>", "<integration-key>")
				c.ConsumesCommandType(fixtures.MessageC{})
				c.ProducesEventType(fixtures.MessageF{})
			},
		}

		projection = &fixtures.ProjectionMessageHandler{
			ConfigureFunc: func(c dogma.ProjectionConfigurer) {
				c.Identity("<projection>", "<projection-key>")
				c.ConsumesEventType(fixtures.MessageD{})
				c.ConsumesEventType(fixtures.MessageE{}) // shared with <process>
			},
		}

		app = &fixtures.Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "<app-key>")
				c.RegisterAggregate(aggregate)
				c.RegisterProcess(process)
				c.RegisterIntegration(integration)
				c.RegisterProjection(projection)
			},
		}
	})

	When("the configuration is valid", func() {
		var cfg *RichApplication

		BeforeEach(func() {
			var err error
			cfg, err = NewApplicationConfig(app)
			Expect(err).ShouldNot(HaveOccurred())
		})

		Describe("func Identity()", func() {

		})

		Describe("func TypeName()", func() {

		})

		Describe("func Messages()", func() {
			It("provides complete role map", func() {
				Expect(cfg.Messages()).To(Equal(
					map[TypeName]MessageRole{
						configfixtures.MessageATypeName: CommandMessageRole,
						configfixtures.MessageBTypeName: EventMessageRole,
						configfixtures.MessageCTypeName: CommandMessageRole,
						configfixtures.MessageDTypeName: EventMessageRole,
						configfixtures.MessageETypeName: EventMessageRole,
						configfixtures.MessageFTypeName: EventMessageRole,
						configfixtures.MessageTTypeName: TimeoutMessageRole,
					},
				))
			})
		})

		Describe("func ConsumedMessages()", func() {

		})

		Describe("func ProducedMessages()", func() {

		})

		Describe("func AcceptVisitor()", func() {

		})

		Describe("func ReflectType()", func() {

		})

		Describe("func ReflectTypeOf()", func() {

		})

		Describe("func MessageTypeOf()", func() {

		})

		Describe("func AcceptRichVisitor()", func() {

		})

		Describe("func Handlers()", func() {

		})

		Describe("func HandlerByIdentity()", func() {

		})

		Describe("func HandlerByName()", func() {

		})

		Describe("func HandlerByKey()", func() {

		})

		Describe("func ConsumersOf()", func() {

		})

		Describe("func ProducersOf()", func() {

		})

		Describe("func ForeignMessages()", func() {

		})

		Describe("func RichHandlers()", func() {

		})

		Describe("func RichHandlerByIdentity()", func() {

		})

		Describe("func RichHandlerByName()", func() {

		})

		Describe("func RichHandlerByKey()", func() {

		})

		Describe("func RichConsumersOf()", func() {

		})

		Describe("func RichProducersOf()", func() {

		})

		It("sets the consumers map", func() {
			Expect(cfg.Consumers).To(Equal(
				map[MessageType][]HandlerConfig{
					MessageAType: {cfg.HandlersByName["<aggregate>"]},
					MessageBType: {cfg.HandlersByName["<process>"]},
					MessageCType: {cfg.HandlersByName["<integration>"]},
					MessageDType: {cfg.HandlersByName["<projection>"]},
					MessageEType: {cfg.HandlersByName["<process>"], cfg.HandlersByName["<projection>"]},
					MessageTType: {cfg.HandlersByName["<process>"]},
				},
			))
		})

		It("sets the producers map", func() {
			Expect(cfg.Producers).To(Equal(
				map[MessageType][]HandlerConfig{
					MessageCType: {cfg.HandlersByName["<process>"]},
					MessageEType: {cfg.HandlersByName["<aggregate>"]},
					MessageFType: {cfg.HandlersByName["<integration>"]},
					MessageTType: {cfg.HandlersByName["<process>"]},
				},
			))
		})

		Describe("func Identity()", func() {
			It("returns the app identity", func() {
				Expect(cfg.Identity()).To(Equal(
					MustNewIdentity("<app>", "<app-key>"),
				))
			})
		})
	})

	When("the app does not configure an identity", func() {
		BeforeEach(func() {
			app.ConfigureFunc = nil
		})

		It("returns a descriptive error", func() {
			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.Application.Configure() did not call ApplicationConfigurer.Identity()`,
				),
			))
		})
	})

	When("the app configures multiple identities", func() {
		BeforeEach(func() {
			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
				c.Identity("<name>", "<key>")
				c.Identity("<other>", "<key>")
			}
		})

		It("returns a descriptive error", func() {
			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.Application.Configure() has already called ApplicationConfigurer.Identity("<name>", "<key>")`,
				),
			))
		})
	})

	When("the app configures an invalid application name", func() {
		BeforeEach(func() {
			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
				c.Identity("\t \n", "<app-key>")
			}
		})

		It("returns a descriptive error", func() {
			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid name "\t \n", names must be non-empty, printable UTF-8 strings with no whitespace`,
				),
			))
		})
	})

	When("the app configures an invalid application key", func() {
		BeforeEach(func() {
			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "\t \n")
			}
		})

		It("returns a descriptive error", func() {
			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.Application.Configure() called ApplicationConfigurer.Identity() with an invalid key "\t \n", keys must be non-empty, printable UTF-8 strings with no whitespace`,
				),
			))
		})
	})

	When("the app contains an invalid handler configurations", func() {
		It("returns an error when an aggregate is misconfigured", func() {
			aggregate.ConfigureFunc = nil

			_, err := NewApplicationConfig(app)

			Expect(err).Should(HaveOccurred())
		})

		It("returns an error when a process is misconfigured", func() {
			process.ConfigureFunc = nil

			_, err := NewApplicationConfig(app)

			Expect(err).Should(HaveOccurred())
		})

		It("returns an error when an integration is misconfigured", func() {
			integration.ConfigureFunc = nil

			_, err := NewApplicationConfig(app)

			Expect(err).Should(HaveOccurred())
		})

		It("returns an error when a projection is misconfigured", func() {
			projection.ConfigureFunc = nil

			_, err := NewApplicationConfig(app)

			Expect(err).Should(HaveOccurred())
		})
	})

	When("the app contains conflicting handler identities", func() {
		It("returns an error when an aggregate name is in conflict", func() {
			aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
				c.Identity("<process>", "<aggregate-key>") // conflict!
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			}

			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "<app-key>")
				c.RegisterProcess(process)
				c.RegisterAggregate(aggregate) // register the conflicting aggregate last
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.AggregateMessageHandler can not use the handler name "<process>", because it is already used by *fixtures.ProcessMessageHandler`,
				),
			))
		})

		It("returns an error when an aggregate key is in conflict", func() {
			aggregate.ConfigureFunc = func(c dogma.AggregateConfigurer) {
				c.Identity("<aggregate>", "<process-key>") // conflict!
				c.ConsumesCommandType(fixtures.MessageA{})
				c.ProducesEventType(fixtures.MessageE{})
			}

			app.ConfigureFunc = func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "<app-key>")
				c.RegisterProcess(process)
				c.RegisterAggregate(aggregate) // register the conflicting aggregate last
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.AggregateMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
				),
			))
		})

		It("returns an error when a process name is in conflict", func() {
			process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
				c.Identity("<aggregate>", "<process-key>") // conflict!
				c.ConsumesEventType(fixtures.MessageB{})
				c.ProducesCommandType(fixtures.MessageC{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.ProcessMessageHandler can not use the handler name "<aggregate>", because it is already used by *fixtures.AggregateMessageHandler`,
				),
			))
		})

		It("returns an error when a process key is in conflict", func() {
			process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
				c.Identity("<process>", "<aggregate-key>") // conflict!
				c.ConsumesEventType(fixtures.MessageB{})
				c.ProducesCommandType(fixtures.MessageC{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.ProcessMessageHandler can not use the handler key "<aggregate-key>", because it is already used by *fixtures.AggregateMessageHandler`,
				),
			))
		})

		It("returns an error when an integration name is in conflict", func() {
			integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
				c.Identity("<process>", "<integration-key>") // conflict!
				c.ConsumesCommandType(fixtures.MessageC{})
				c.ProducesEventType(fixtures.MessageF{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.IntegrationMessageHandler can not use the handler name "<process>", because it is already used by *fixtures.ProcessMessageHandler`,
				),
			))
		})

		It("returns an error when an integration key is in conflict", func() {
			integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
				c.Identity("<integration>", "<process-key>") // conflict!
				c.ConsumesCommandType(fixtures.MessageC{})
				c.ProducesEventType(fixtures.MessageF{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.IntegrationMessageHandler can not use the handler key "<process-key>", because it is already used by *fixtures.ProcessMessageHandler`,
				),
			))
		})

		It("returns an error when a projection name is in conflict", func() {
			projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
				c.Identity("<integration>", "<projection-key>") // conflict!
				c.ConsumesEventType(fixtures.MessageD{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.ProjectionMessageHandler can not use the handler name "<integration>", because it is already used by *fixtures.IntegrationMessageHandler`,
				),
			))
		})

		It("returns an error when a projection key is in conflict", func() {
			projection.ConfigureFunc = func(c dogma.ProjectionConfigurer) {
				c.Identity("<projection>", "<integration-key>") // conflict!
				c.ConsumesEventType(fixtures.MessageD{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`*fixtures.ProjectionMessageHandler can not use the handler key "<integration-key>", because it is already used by *fixtures.IntegrationMessageHandler`,
				),
			))
		})
	})

	It("returns an error when the app contains multiple consumers of the same command", func() {
		integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
			c.Identity("<integration>", "<integration-key>")
			c.ConsumesCommandType(fixtures.MessageA{}) // conflict with <aggregate>
			c.ProducesEventType(fixtures.MessageF{})
		}

		_, err := NewApplicationConfig(app)

		Expect(err).To(Equal(
			ValidationError(
				`the "<integration>" handler can not consume fixtures.MessageA commands because they are already consumed by "<aggregate>"`,
			),
		))
	})

	It("returns an error when the app contains multiple producers of the same event", func() {
		integration.ConfigureFunc = func(c dogma.IntegrationConfigurer) {
			c.Identity("<integration>", "<integration-key>")
			c.ConsumesCommandType(fixtures.MessageC{})
			c.ProducesEventType(fixtures.MessageE{}) // conflict with <aggregate>
		}

		_, err := NewApplicationConfig(app)

		Expect(err).To(Equal(
			ValidationError(
				`the "<integration>" handler can not produce fixtures.MessageE events because they are already produced by "<aggregate>"`,
			),
		))
	})

	It("does not return an error when the app contains multiple processes that schedule the same timeout", func() {
		process1 := &fixtures.ProcessMessageHandler{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("<process-1>", "<process-1-key>")
				c.ConsumesEventType(fixtures.MessageB{})
				c.ProducesCommandType(fixtures.MessageC{})
				c.SchedulesTimeoutType(fixtures.MessageT{})
			},
		}

		process2 := &fixtures.ProcessMessageHandler{
			ConfigureFunc: func(c dogma.ProcessConfigurer) {
				c.Identity("<process-2>", "<process-2-key>")
				c.ConsumesEventType(fixtures.MessageB{})
				c.ProducesCommandType(fixtures.MessageC{})
				c.SchedulesTimeoutType(fixtures.MessageT{})
			},
		}

		app := &Application{
			ConfigureFunc: func(c dogma.ApplicationConfigurer) {
				c.Identity("<app>", "<app-key>")
				c.RegisterProcess(process1)
				c.RegisterProcess(process2)
			},
		}

		_, err := NewApplicationConfig(app)

		Expect(err).ShouldNot(HaveOccurred())
	})

	When("multiple handlers use a single message type in differing roles", func() {
		It("returns an error when a conflict occurs with a consumed message", func() {
			process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
				c.Identity("<process>", "<process-key>")
				c.ConsumesEventType(fixtures.MessageA{}) // conflict with <aggregate>
				c.ProducesCommandType(fixtures.MessageC{})
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`the "<process>" handler configures fixtures.MessageA as an event but "<aggregate>" configures it as a command`,
				),
			))
		})

		It("returns an error when a conflict occurs with a produced message", func() {
			process.ConfigureFunc = func(c dogma.ProcessConfigurer) {
				c.Identity("<process>", "<process-key>")
				c.ConsumesEventType(fixtures.MessageB{})
				c.ProducesCommandType(fixtures.MessageE{}) // conflict with <aggregate>
			}

			_, err := NewApplicationConfig(app)

			Expect(err).To(Equal(
				ValidationError(
					`the "<process>" handler configures fixtures.MessageE as a command but "<aggregate>" configures it as an event`,
				),
			))
		})
	})
})
