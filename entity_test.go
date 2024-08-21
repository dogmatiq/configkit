package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message" // can't dot-import due to conflicts
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EntityMessageNames", func() {
	Describe("func RoleOf()", func() {
		It("returns the role of a produced message", func() {
			m := EntityMessageNames{
				Produced: message.NameRoles{
					message.NameFor[CommandStub[TypeA]](): message.CommandRole,
				},
			}

			r, ok := m.RoleOf(message.NameFor[CommandStub[TypeA]]())
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(message.CommandRole))
		})

		It("returns the role of a consumed message", func() {
			m := EntityMessageNames{
				Consumed: message.NameRoles{
					message.NameFor[CommandStub[TypeA]](): message.CommandRole,
				},
			}

			r, ok := m.RoleOf(message.NameFor[CommandStub[TypeA]]())
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(message.CommandRole))
		})

		It("returns false if the message is neither produced nor consumed", func() {
			m := EntityMessageNames{}

			_, ok := m.RoleOf(message.NameFor[CommandStub[TypeA]]())
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func All()", func() {
		It("returns the union of the produced and consumed messages", func() {
			m := EntityMessageNames{
				Produced: message.NameRoles{
					message.NameFor[CommandStub[TypeA]](): message.CommandRole,
				},
				Consumed: message.NameRoles{
					message.NameFor[EventStub[TypeA]](): message.EventRole,
				},
			}

			Expect(m.All()).To(Equal(
				message.NameRoles{
					message.NameFor[CommandStub[TypeA]](): message.CommandRole,
					message.NameFor[EventStub[TypeA]]():   message.EventRole,
				},
			))
		})
	})

	Describe("func ForeignMessageNames()", func() {
		It("returns the set of messages that belong to another application", func() {
			m := EntityMessageNames{
				Produced: message.NameRoles{
					message.NameFor[EventStub[TypeA]]():   message.EventRole,
					message.NameFor[CommandStub[TypeA]](): message.CommandRole, // foreign-consumed command
				},
				Consumed: message.NameRoles{
					message.NameFor[CommandStub[TypeB]](): message.CommandRole, // foreign-produced command
					message.NameFor[EventStub[TypeB]]():   message.EventRole,   // foreign-produced event
					message.NameFor[EventStub[TypeA]]():   message.EventRole,
				},
			}

			Expect(m.Foreign()).To(Equal(
				EntityMessageNames{
					Produced: message.NameRoles{
						message.NameFor[CommandStub[TypeA]](): message.CommandRole,
					},
					Consumed: message.NameRoles{
						message.NameFor[CommandStub[TypeB]](): message.CommandRole,
						message.NameFor[EventStub[TypeB]]():   message.EventRole,
					},
				},
			))
		})
	})

	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageNames{
				Produced: message.NameRoles{
					message.NameFor[EventStub[TypeA]](): message.EventRole,
				},
				Consumed: message.NameRoles{
					message.NameFor[CommandStub[TypeA]](): message.CommandRole,
				},
			}

			b := EntityMessageNames{
				Produced: message.NameRoles{
					message.NameFor[EventStub[TypeA]](): message.EventRole,
				},
				Consumed: message.NameRoles{
					message.NameFor[CommandStub[TypeA]](): message.CommandRole,
				},
			}

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessageNames) {
				a := EntityMessageNames{
					Produced: message.NameRoles{
						message.NameFor[EventStub[TypeA]](): message.EventRole,
					},
					Consumed: message.NameRoles{
						message.NameFor[CommandStub[TypeA]](): message.CommandRole,
					},
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"produced messages differ",
				EntityMessageNames{
					Produced: message.NameRoles{
						message.NameFor[EventStub[TypeA]]():   message.EventRole,
						message.NameFor[TimeoutStub[TypeA]](): message.TimeoutRole, // diff
					},
					Consumed: message.NameRoles{
						message.NameFor[CommandStub[TypeA]](): message.CommandRole,
					},
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessageNames{
					Produced: message.NameRoles{
						message.NameFor[EventStub[TypeA]](): message.EventRole,
					},
					Consumed: message.NameRoles{
						message.NameFor[CommandStub[TypeA]](): message.CommandRole,
						message.NameFor[TimeoutStub[TypeA]](): message.TimeoutRole, // diff
					},
				},
			),
		)
	})
})

var _ = Describe("type EntityMessageTypes", func() {
	Describe("func RoleOf()", func() {
		It("returns the role of a produced message", func() {
			m := EntityMessageTypes{
				Produced: message.TypeRoles{
					message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
				},
			}

			r, ok := m.RoleOf(message.TypeFor[CommandStub[TypeA]]())
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(message.CommandRole))
		})

		It("returns the role of a consumed message", func() {
			m := EntityMessageTypes{
				Consumed: message.TypeRoles{
					message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
				},
			}

			r, ok := m.RoleOf(message.TypeFor[CommandStub[TypeA]]())
			Expect(ok).To(BeTrue())
			Expect(r).To(Equal(message.CommandRole))
		})

		It("returns false if the message is neither produced nor consumed", func() {
			m := EntityMessageTypes{}

			_, ok := m.RoleOf(message.TypeFor[CommandStub[TypeA]]())
			Expect(ok).To(BeFalse())
		})
	})

	Describe("func All()", func() {
		It("returns the union of the produced and consumed messages", func() {
			m := EntityMessageTypes{
				Produced: message.TypeRoles{
					message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
				},
				Consumed: message.TypeRoles{
					message.TypeFor[EventStub[TypeA]](): message.EventRole,
				},
			}

			Expect(m.All()).To(Equal(
				message.TypeRoles{
					message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
					message.TypeFor[EventStub[TypeA]]():   message.EventRole,
				},
			))
		})
	})

	Describe("func Foreign()", func() {
		It("returns the set of messages that belong to another entity", func() {
			m := EntityMessageTypes{
				Produced: message.TypeRoles{
					message.TypeFor[EventStub[TypeA]]():   message.EventRole,
					message.TypeFor[CommandStub[TypeA]](): message.CommandRole, // foreign-consumed command
				},
				Consumed: message.TypeRoles{
					message.TypeFor[CommandStub[TypeB]](): message.CommandRole, // foreign-produced command
					message.TypeFor[EventStub[TypeB]]():   message.EventRole,   // foreign-produced event
					message.TypeFor[EventStub[TypeA]]():   message.EventRole,
				},
			}

			Expect(m.Foreign()).To(Equal(
				EntityMessageTypes{
					Produced: message.TypeRoles{
						message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
					},
					Consumed: message.TypeRoles{
						message.TypeFor[CommandStub[TypeB]](): message.CommandRole,
						message.TypeFor[EventStub[TypeB]]():   message.EventRole,
					},
				},
			))
		})
	})

	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageTypes{
				Produced: message.TypeRoles{
					message.TypeFor[EventStub[TypeA]](): message.EventRole,
				},
				Consumed: message.TypeRoles{
					message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
				},
			}

			b := EntityMessageTypes{
				Produced: message.TypeRoles{
					message.TypeFor[EventStub[TypeA]](): message.EventRole,
				},
				Consumed: message.TypeRoles{
					message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
				},
			}

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessageTypes) {
				a := EntityMessageTypes{
					Produced: message.TypeRoles{
						message.TypeFor[EventStub[TypeA]](): message.EventRole,
					},
					Consumed: message.TypeRoles{
						message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
					},
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"produced messages differ",
				EntityMessageTypes{
					Produced: message.TypeRoles{
						message.TypeFor[EventStub[TypeA]]():   message.EventRole,
						message.TypeFor[TimeoutStub[TypeA]](): message.TimeoutRole, // diff
					},
					Consumed: message.TypeRoles{
						message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
					},
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessageTypes{
					Produced: message.TypeRoles{
						message.TypeFor[EventStub[TypeA]](): message.EventRole,
					},
					Consumed: message.TypeRoles{
						message.TypeFor[CommandStub[TypeA]](): message.CommandRole,
						message.TypeFor[TimeoutStub[TypeA]](): message.TimeoutRole, // diff
					},
				},
			),
		)
	})
})
