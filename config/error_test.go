package config_test

import (
	. "github.com/dogmatiq/configkit/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("type Error", func() {
	Describe("func Error", func() {
		It("returns the error message", func() {
			err := Error("<message>")
			Expect(err.Error()).To(Equal("<message>"))
		})
	})
})

var _ = Describe("func Recover()", func() {
	It("recovers from config related panics", func() {
		err := func() (err error) {
			defer Recover(&err)
			Panicf("<value>")
			return nil
		}()

		Expect(err).To(Equal(Error("<value>")))
	})

	It("does not recover from unrelated panics", func() {
		var value interface{}

		func() {
			defer func() {
				value = recover()
			}()

			func() (err error) {
				defer Recover(&err)
				panic("<value>") // not a configkit.Error
			}()
		}()

		Expect(value).To(Equal("<value>"))
	})

	It("does not panic when no panic occurs", func() {
		err := func() (err error) {
			defer Recover(&err)
			return nil
		}()

		Expect(err).ShouldNot(HaveOccurred())
	})

	It("panics when passed a nil pointer", func() {
		Expect(func() {
			Recover(nil)
		}).To(Panic())
	})
})
