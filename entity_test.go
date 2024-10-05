package configkit_test

import (
	"maps"

	. "github.com/dogmatiq/configkit"
	. "github.com/dogmatiq/enginekit/enginetest/stubs"
	"github.com/dogmatiq/enginekit/message"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("type EntityMessages", func() {
	Describe("func IsEqual()", func() {
		It("returns true if the sets are equivalent", func() {
			a := EntityMessages[message.Name]{
				message.NameOf(EventA1): {
					Kind:       message.EventKind,
					IsProduced: true,
				},
				message.NameOf(CommandA1): {
					Kind:       message.CommandKind,
					IsConsumed: true,
				},
			}

			b := maps.Clone(a)

			Expect(a.IsEqual(b)).To(BeTrue())
		})

		DescribeTable(
			"returns false if the sets are not equivalent",
			func(b EntityMessages[message.Name]) {
				a := EntityMessages[message.Name]{
					message.NameOf(EventA1): {
						Kind:       message.EventKind,
						IsProduced: true,
					},
					message.NameOf(CommandA1): {
						Kind:       message.CommandKind,
						IsConsumed: true,
					},
				}
				Expect(a.IsEqual(b)).To(BeFalse())
			},
			Entry(
				"kinds differ",
				EntityMessages[message.Name]{
					message.NameOf(EventA1): {
						Kind:       message.EventKind,
						IsProduced: true,
					},
					message.NameOf(CommandA1): {
						Kind:       message.TimeoutKind, // diff
						IsConsumed: true,
					},
				},
			),
			Entry(
				"produced messages differ",
				EntityMessages[message.Name]{
					message.NameOf(EventA1): {
						Kind:       message.EventKind,
						IsProduced: true,
					},
					message.NameOf(CommandA1): {
						Kind:       message.CommandKind,
						IsConsumed: true,
						IsProduced: true, // diff
					},
				},
			),
			Entry(
				"consumed messages differ",
				EntityMessages[message.Name]{
					message.NameOf(EventA1): {
						Kind:       message.EventKind,
						IsProduced: true,
						IsConsumed: true, // diff
					},
					message.NameOf(CommandA1): {
						Kind:       message.CommandKind,
						IsConsumed: true,
					},
				},
			),
			Entry(
				"keys differ",
				EntityMessages[message.Name]{
					message.NameOf(EventB1): { // diff
						Kind:       message.EventKind,
						IsProduced: true,
					},
					message.NameOf(CommandA1): {
						Kind:       message.CommandKind,
						IsConsumed: true,
					},
				},
			),
			Entry(
				"lengths differ",
				EntityMessages[message.Name]{
					message.NameOf(EventA1): {
						Kind:       message.EventKind,
						IsProduced: true,
					},
					message.NameOf(CommandA1): {
						Kind:       message.CommandKind,
						IsConsumed: true,
					},
					// diff
					message.NameOf(TimeoutA1): {
						Kind:       message.TimeoutKind,
						IsConsumed: true,
						IsProduced: true,
					},
				},
			),
		)
	})
})
