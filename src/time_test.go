package Time_test

import (
	"time"

	Time "github.com/go-composites/time/src"
	Duration "github.com/go-composites/time/src/duration"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Time", func() {

	ginkgo.Describe("constructors", func() {
		ginkgo.It("builds from a Unix timestamp", func() {
			t := Time.FromUnix(0)
			gomega.Expect(t.ToUnix()).To(gomega.BeEquivalentTo(0))
			gomega.Expect(t.IsNull()).To(gomega.BeFalse())
		})

		ginkgo.Describe("Parse", func() {
			ginkgo.It("parses a well-formed value", func() {
				r := Time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
				gomega.Expect(r.HasError()).To(gomega.BeFalse())
				gomega.Expect(r.Payload().(Time.Interface).ToUnix()).To(gomega.BeEquivalentTo(0))
			})
			ginkgo.It("returns a Result carrying an error on a malformed value", func() {
				r := Time.Parse(time.RFC3339, "not-a-time")
				gomega.Expect(r.HasError()).To(gomega.BeTrue())
				gomega.Expect(r.Error().Message()).NotTo(gomega.BeEmpty())
			})
		})
	})

	ginkgo.Describe("conversions", func() {
		ginkgo.It("formats according to a layout", func() {
			gomega.Expect(Time.FromUnix(0).Format(time.RFC3339)).To(gomega.Equal("1970-01-01T00:00:00Z"))
		})
		ginkgo.It("renders RFC3339 from ToGoString", func() {
			gomega.Expect(Time.FromUnix(0).ToGoString()).To(gomega.Equal("1970-01-01T00:00:00Z"))
		})
	})

	ginkgo.Describe("comparisons", func() {
		var early = Time.FromUnix(0)
		var late = Time.FromUnix(100)

		ginkgo.It("reports before", func() {
			gomega.Expect(early.Before(late)).To(gomega.BeTrue())
			gomega.Expect(late.Before(early)).To(gomega.BeFalse())
		})
		ginkgo.It("reports after", func() {
			gomega.Expect(late.After(early)).To(gomega.BeTrue())
			gomega.Expect(early.After(late)).To(gomega.BeFalse())
		})
		ginkgo.It("reports equality", func() {
			gomega.Expect(early.Equal(Time.FromUnix(0))).To(gomega.BeTrue())
			gomega.Expect(early.Equal(late)).To(gomega.BeFalse())
		})
	})

	ginkgo.Describe("arithmetic", func() {
		ginkgo.It("Add shifts a Time forward by a Duration", func() {
			r := Time.FromUnix(0).Add(Duration.FromSeconds(100))
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(Time.Interface).ToUnix()).To(gomega.BeEquivalentTo(100))
		})
		ginkgo.It("Sub yields the Duration between two Times", func() {
			gap := Time.FromUnix(100).Sub(Time.FromUnix(40))
			gomega.Expect(gap.ToSeconds()).To(gomega.BeEquivalentTo(60))
		})
	})

	ginkgo.Describe("the Null-Object variant", func() {
		var n = Time.Null()

		ginkgo.It("satisfies the Time interface", func() {
			var _ Time.Interface = n
		})
		ginkgo.It("reports IsNull() true", func() {
			gomega.Expect(n.IsNull()).To(gomega.BeTrue())
		})
		ginkgo.It("denotes the zero instant", func() {
			gomega.Expect(n.ToUnix()).To(gomega.BeEquivalentTo(0))
		})
		ginkgo.It("formats to the empty string", func() {
			gomega.Expect(n.Format(time.RFC3339)).To(gomega.Equal(``))
		})
		ginkgo.It("renders the null marker", func() {
			gomega.Expect(n.ToGoString()).To(gomega.Equal(`<NullTime>`))
		})
		ginkgo.It("Before is always false", func() {
			gomega.Expect(n.Before(Time.FromUnix(0))).To(gomega.BeFalse())
		})
		ginkgo.It("After is always false", func() {
			gomega.Expect(n.After(Time.FromUnix(0))).To(gomega.BeFalse())
		})
		ginkgo.It("Equal is true only against another null", func() {
			gomega.Expect(n.Equal(Time.Null())).To(gomega.BeTrue())
			gomega.Expect(n.Equal(Time.FromUnix(0))).To(gomega.BeFalse())
		})
		ginkgo.It("Add returns a Result whose payload is a null Time", func() {
			r := n.Add(Duration.FromSeconds(60))
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			gomega.Expect(r.Payload().(Time.Interface).IsNull()).To(gomega.BeTrue())
		})
		ginkgo.It("Sub returns a null Duration", func() {
			gomega.Expect(n.Sub(Time.FromUnix(0)).IsNull()).To(gomega.BeTrue())
		})
	})
})
