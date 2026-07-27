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

		ginkgo.Describe("ParseAny", func() {
			ginkgo.It("parses layout-free real-world dates (2-digit year + named zone)", func() {
				r := Time.ParseAny("Tue, 07 Jul 26 11:13:37 UTC")
				gomega.Expect(r.HasError()).To(gomega.BeFalse())
				got := r.Payload().(Time.Interface)
				gomega.Expect(got.Format("2006-01-02")).To(gomega.Equal("2026-07-07"))
			})
			ginkgo.It("returns a Result carrying an error on an unrecognisable value", func() {
				r := Time.ParseAny("not-a-time")
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

	ginkgo.Describe("time zones", func() {
		ginkgo.It("relocates an instant into a known IANA zone", func() {
			r := Time.FromUnix(0).In("Europe/Paris")
			gomega.Expect(r.HasError()).To(gomega.BeFalse())
			relocated := r.Payload().(Time.Interface)
			gomega.Expect(relocated.Zone()).To(gomega.Equal("CET"))
		})
		ginkgo.It("returns a Result carrying an error on an unknown zone", func() {
			r := Time.FromUnix(0).In("Mars/Olympus_Mons")
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).NotTo(gomega.BeEmpty())
		})
		ginkgo.It("reports the zone abbreviation of a UTC instant", func() {
			gomega.Expect(Time.FromUnix(0).Zone()).To(gomega.Equal("UTC"))
		})
		ginkgo.It("UTC sets the location to UTC", func() {
			r := Time.FromUnix(0).In("Europe/Paris")
			relocated := r.Payload().(Time.Interface)
			gomega.Expect(relocated.UTC().Zone()).To(gomega.Equal("UTC"))
		})
		ginkgo.It("preserves the instant across In and UTC (equality is unchanged)", func() {
			epoch := Time.FromUnix(0)
			paris := epoch.In("Europe/Paris").Payload().(Time.Interface)
			gomega.Expect(epoch.Equal(paris)).To(gomega.BeTrue())
			gomega.Expect(paris.Equal(epoch)).To(gomega.BeTrue())
			gomega.Expect(epoch.Equal(paris.UTC())).To(gomega.BeTrue())
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
		ginkgo.It("In returns a Result carrying an error", func() {
			r := n.In("Europe/Paris")
			gomega.Expect(r.HasError()).To(gomega.BeTrue())
			gomega.Expect(r.Error().Message()).NotTo(gomega.BeEmpty())
		})
		ginkgo.It("Zone is the empty string", func() {
			gomega.Expect(n.Zone()).To(gomega.Equal(``))
		})
		ginkgo.It("UTC returns a null Time", func() {
			gomega.Expect(n.UTC().IsNull()).To(gomega.BeTrue())
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
