package static_test

import (
	"github.com/dogmatiq/configkit"
	. "github.com/onsi/gomega"
)

// matchIdentities matches the given identities to those of the handlers in the
// handler set.
func matchIdentities(
	hs configkit.HandlerSet,
	identities ...configkit.Identity,
) {
	// TODO: should this be exhaustive?

	for _, identity := range identities {
		_, ok := hs.ByIdentity(identity)
		Expect(ok).To(BeTrue())
	}
}
