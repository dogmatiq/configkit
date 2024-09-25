package configkit_test

import (
	. "github.com/dogmatiq/configkit"
	"github.com/dogmatiq/configkit/message"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EntityMessageNames", func() {
	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageNames{
				Kinds: map[message.Name]message.Kind{
					message.NameOf(EventA1):   message.EventKind,
					message.NameOf(CommandA1): message.CommandKind,
				},
				Produced: message.NamesOf(EventA1),
				Consumed: message.NamesOf(CommandA1),
			}

			b := EntityMessageNames{
				Kinds: map[message.Name]message.Kind{
					message.NameOf(EventA1):   message.EventKind,
					message.NameOf(CommandA1): message.CommandKind,
				},
				Produced: message.NamesOf(EventA1),
				Consumed: message.NamesOf(CommandA1),
			}

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessageNames) {
				a := EntityMessageNames{
					Kinds: map[message.Name]message.Kind{
						message.NameOf(EventA1):   message.EventKind,
						message.NameOf(CommandA1): message.CommandKind,
					},
					Produced: message.NamesOf(EventA1),
					Consumed: message.NamesOf(CommandA1),
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"kinds differ",
				EntityMessageNames{
					Kinds: map[message.Name]message.Kind{
						message.NameOf(EventA1):   message.EventKind,
						message.NameOf(CommandA1): message.TimeoutKind, // diff
					},
					Produced: message.NamesOf(EventA1),
					Consumed: message.NamesOf(CommandA1),
				},
			),
			Entry(
				"produced messages differ",
				EntityMessageNames{
					Kinds: map[message.Name]message.Kind{
						message.NameOf(EventA1):   message.EventKind,
						message.NameOf(CommandA1): message.CommandKind,
					},
					Produced: message.NamesOf(
						EventA1,
						CommandA1, // diff
					),
					Consumed: message.NamesOf(CommandA1),
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessageNames{
					Kinds: map[message.Name]message.Kind{
						message.NameOf(EventA1):   message.EventKind,
						message.NameOf(CommandA1): message.CommandKind,
					},
					Produced: message.NamesOf(EventA1),
					Consumed: message.NamesOf(
						EventA1,
						CommandA1, // diff
					),
				},
			),
		)
	})
})

var _ = Describe("type EntityMessageTypes", func() {
	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessageTypes{
				Produced: message.TypesOf(EventA1),
				Consumed: message.TypesOf(CommandA1),
			}

			b := EntityMessageTypes{
				Produced: message.TypesOf(EventA1),
				Consumed: message.TypesOf(CommandA1),
			}

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessageTypes) {
				a := EntityMessageTypes{
					Produced: message.TypesOf(EventA1),
					Consumed: message.TypesOf(CommandA1),
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"produced messages differ",
				EntityMessageTypes{
					Produced: message.TypesOf(
						EventA1,
						TimeoutA1, // diff
					),
					Consumed: message.TypesOf(CommandA1),
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessageTypes{
					Produced: message.TypesOf(EventA1),
					Consumed: message.TypesOf(
						CommandA1,
						TimeoutA1, // diff
					),
				},
			),
		)
	})
})
