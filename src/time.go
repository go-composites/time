package Time

import (
	"time"
	// Embedding the IANA time-zone database makes time.LoadLocation work on
	// every architecture and inside the toolchain-less CI containers, without
	// relying on a system tzdata install — keeping In() deterministic.
	_ "time/tzdata"

	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
	Duration "github.com/go-composites/time/src/duration"

	dates "github.com/go-datetime/dates"
)

/*
Time is a composite over an instant in time, wrapping Go's stdlib time.Time
(Ruby's Time as the reference).

It follows Composition-Oriented Programming: it is interface-first, its fallible
constructors (Parse, Add) return a Result.Interface so that failures are values
rather than panics, and it ships a Null-Object variant so callers never have to
test for a bare nil.

Time is deterministic by construction: there is no Now(). Instances are built
only from explicit values (FromUnix, Parse), which keeps behaviour reproducible
across architectures and test runs.
*/
type Interface interface {
	ToUnix() int64
	Format(layout string) string
	ToGoString() string
	In(name string) Result.Interface
	Zone() string
	UTC() Interface
	Before(Interface) bool
	After(Interface) bool
	Equal(Interface) bool
	Add(Duration.Interface) Result.Interface
	Sub(Interface) Duration.Interface
	IsNull() bool
}

type data struct {
	value time.Time
}

/*
FromUnix is the Time constructor from a Unix timestamp (seconds since the
epoch), in UTC.

	t := Time.FromUnix(0) // 1970-01-01T00:00:00Z
*/
func FromUnix(sec int64) Interface {
	return &data{value: time.Unix(sec, 0).UTC()}
}

/*
fromGo wraps a Go time.Time into a Time composite.
*/
func fromGo(value time.Time) Interface {
	return &data{value: value}
}

/*
Parse builds a Time from a layout and a value (see time.Parse).

On success the Result carries the Time as its payload; when the value does not
match the layout the Result carries an Error instead — the operation never
panics and never returns nil.

	r := Time.Parse(time.RFC3339, "2026-06-20T12:00:00Z")
	if r.HasError() {
	    // r.Error().Message()
	}
*/
func Parse(layout, value string) Result.Interface {
	parsed, err := time.Parse(layout, value)
	if err != nil {
		return Result.New(
			Result.WithError(
				Error.New(err.Error()),
			),
		)
	}
	return Result.New(
		Result.WithPayload(fromGo(parsed)),
	)
}

/*
ParseAny builds a Time from a date/time string WITHOUT an explicit layout,
trying the real-world format zoo — RFC 822/1123/2822/3339/5322, asctime, HTTP,
2-digit years, ISO 8601 and named-zone abbreviations — via the shared
go-datetime/dates parser. This is the lenient, MRI-style Time.parse behaviour
(Ruby's Time.parse accepts almost any recognisable date), so consumers such as
the rbgo interpreter delegate here instead of reimplementing a layout table.

On success the Result carries the Time; an unrecognisable value carries an
Error; the operation never panics and never returns nil.

	r := Time.ParseAny("Tue, 07 Jul 26 11:13:37 UTC") // 2-digit year + named zone
*/
func ParseAny(value string) Result.Interface {
	parsed, err := dates.Parse(value)
	if err != nil {
		return Result.New(
			Result.WithError(
				Error.New(err.Error()),
			),
		)
	}
	return Result.New(
		Result.WithPayload(fromGo(parsed)),
	)
}

/*
ToUnix returns the instant as a Unix timestamp (seconds since the epoch).
*/
func (d data) ToUnix() int64 {
	return d.value.Unix()
}

/*
Format renders the instant according to the given layout (see time.Time.Format).
*/
func (d data) Format(layout string) string {
	return d.value.Format(layout)
}

/*
ToGoString returns the RFC3339 representation of the instant.
*/
func (d data) ToGoString() string {
	return d.value.Format(time.RFC3339)
}

/*
In returns a Result whose payload is a new Time denoting the SAME instant in
the IANA location name (e.g. "Europe/Paris").

On success the Result carries the relocated Time; when name is not a known zone
the Result carries an Error instead — the operation never panics and never
returns nil. Because only the location changes, the underlying instant (and so
Before/After/Equal) is preserved.

	r := Time.FromUnix(0).In("Europe/Paris")
	if r.HasError() {
	    // r.Error().Message()
	}
*/
func (d data) In(name string) Result.Interface {
	location, err := time.LoadLocation(name)
	if err != nil {
		return Result.New(
			Result.WithError(
				Error.New(err.Error()),
			),
		)
	}
	return Result.New(
		Result.WithPayload(fromGo(d.value.In(location))),
	)
}

/*
Zone returns the abbreviated name of the time zone in effect at the instant
(e.g. "UTC", "CET").
*/
func (d data) Zone() string {
	name, _ := d.value.Zone()
	return name
}

/*
UTC returns a new Time denoting the same instant with its location set to UTC.
*/
func (d data) UTC() Interface {
	return fromGo(d.value.UTC())
}

/*
Before reports whether the receiver is strictly before other.
*/
func (d data) Before(other Interface) bool {
	return d.value.Unix() < other.ToUnix()
}

/*
After reports whether the receiver is strictly after other.
*/
func (d data) After(other Interface) bool {
	return d.value.Unix() > other.ToUnix()
}

/*
Equal reports whether the receiver and other denote the same instant.
*/
func (d data) Equal(other Interface) bool {
	return d.value.Unix() == other.ToUnix()
}

/*
Add returns a Result whose payload is a new Time shifted forward by the given
Duration.
*/
func (d data) Add(duration Duration.Interface) Result.Interface {
	shifted := d.value.Add(time.Duration(duration.ToSeconds()) * time.Second)
	return Result.New(
		Result.WithPayload(fromGo(shifted)),
	)
}

/*
Sub returns the Duration spanning from other to the receiver (the gap between
the two instants).
*/
func (d data) Sub(other Interface) Duration.Interface {
	return Duration.FromSeconds(d.value.Unix() - other.ToUnix())
}

/*
IsNull reports whether the Time is the Null-Object variant.

A concrete Time is never null.
*/
func (d data) IsNull() bool {
	return false
}

type null struct{}

/*
Null returns the Null-Object variant of Time.

It satisfies Interface so callers never have to test for a bare nil: it denotes
the zero instant, its comparisons are false, its Add yields itself, its Sub
yields a null Duration, and IsNull() returns true.
*/
func Null() Interface {
	return &null{}
}

func (n null) ToUnix() int64 {
	return 0
}

func (n null) Format(string) string {
	return ``
}

func (n null) ToGoString() string {
	return `<NullTime>`
}

func (n null) In(string) Result.Interface {
	return Result.New(
		Result.WithError(
			Error.New(`cannot relocate a null Time`),
		),
	)
}

func (n null) Zone() string {
	return ``
}

func (n null) UTC() Interface {
	return Null()
}

func (n null) Before(Interface) bool {
	return false
}

func (n null) After(Interface) bool {
	return false
}

func (n null) Equal(other Interface) bool {
	return other.IsNull()
}

func (n null) Add(Duration.Interface) Result.Interface {
	return Result.New(
		Result.WithPayload(Null()),
	)
}

func (n null) Sub(Interface) Duration.Interface {
	return Duration.Null()
}

func (n null) IsNull() bool {
	return true
}
