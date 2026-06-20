package Duration_test

import (
	Duration "github.com/go-composites/time/src/duration"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Duration", func() {

	ginkgo.Describe("constructors", func() {
		ginkgo.It("builds from a whole number of seconds", func() {
			d := Duration.FromSeconds(90)
			gomega.Expect(d.ToSeconds()).To(gomega.BeEquivalentTo(90))
			gomega.Expect(d.IsNull()).To(gomega.BeFalse())
		})

		ginkgo.Describe("Parse", func() {
			ginkgo.It("parses a well-formed Go duration string", func() {
				r := Duration.Parse("1h30m")
				gomega.Expect(r.HasError()).To(gomega.BeFalse())
				gomega.Expect(r.Payload().(Duration.Interface).ToSeconds()).
					To(gomega.BeEquivalentTo(5400))
			})
			ginkgo.It("returns a Result carrying an error on a malformed string", func() {
				r := Duration.Parse("not-a-duration")
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
				gomega.Expect(r.Error().Message()).NotTo(gomega.BeEmpty())
			})
		})
	})

	ginkgo.Describe("conversions", func() {
		ginkgo.It("renders the Go string representation", func() {
			gomega.Expect(Duration.FromSeconds(5400).ToGoString()).To(gomega.Equal("1h30m0s"))
		})
	})

	ginkgo.Describe("arithmetic", func() {
		var one = Duration.FromSeconds(60)
		var two = Duration.FromSeconds(120)

		ginkgo.It("adds two durations", func() {
			gomega.Expect(one.Add(two).ToSeconds()).To(gomega.BeEquivalentTo(180))
		})
		ginkgo.It("subtracts two durations", func() {
			gomega.Expect(two.Sub(one).ToSeconds()).To(gomega.BeEquivalentTo(60))
		})
	})

	ginkgo.Describe("comparisons", func() {
		ginkgo.It("reports equality", func() {
			gomega.Expect(Duration.FromSeconds(60).Equal(Duration.FromSeconds(60))).To(gomega.BeTrue())
			gomega.Expect(Duration.FromSeconds(60).Equal(Duration.FromSeconds(61))).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("the Null-Object variant", func() {
		var n = Duration.Null()

		ginkgo.It("satisfies the Duration interface", func() {
			var _ Duration.Interface = n
		})
		ginkgo.It("reports IsNull() true", func() {
			gomega.Expect(n.IsNull()).To(gomega.BeTrue())
		})
		ginkgo.It("spans zero seconds", func() {
			gomega.Expect(n.ToSeconds()).To(gomega.BeEquivalentTo(0))
		})
		ginkgo.It("renders the null marker", func() {
			gomega.Expect(n.ToGoString()).To(gomega.Equal(`<NullDuration>`))
		})
		ginkgo.It("Add returns a null duration", func() {
			gomega.Expect(n.Add(Duration.FromSeconds(60)).IsNull()).To(gomega.BeTrue())
		})
		ginkgo.It("Sub returns a null duration", func() {
			gomega.Expect(n.Sub(Duration.FromSeconds(60)).IsNull()).To(gomega.BeTrue())
		})
		ginkgo.It("Equal is true only against another null", func() {
			gomega.Expect(n.Equal(Duration.Null())).To(gomega.BeTrue())
			gomega.Expect(n.Equal(Duration.FromSeconds(60))).To(gomega.BeFalse())
		})
	})
})
