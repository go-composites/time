package Duration

import (
	"time"

	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
)

/*
Duration is a composite over a span of time, wrapping Go's stdlib
time.Duration.

It follows Composition-Oriented Programming: it is interface-first, its
fallible constructor (Parse) returns a Result.Interface so that a malformed
input is a value rather than a panic, and it ships a Null-Object variant so
callers never have to test for a bare nil.
*/
type Interface interface {
	ToSeconds() int64
	ToGoString() string
	Add(Interface) Interface
	Sub(Interface) Interface
	Equal(Interface) bool
	IsNull() bool
}

type data struct {
	value time.Duration
}

/*
FromSeconds is the Duration constructor from a whole number of seconds.

	d := Duration.FromSeconds(90) // 1m30s
*/
func FromSeconds(sec int64) Interface {
	return &data{value: time.Duration(sec) * time.Second}
}

/*
fromGo wraps a Go time.Duration into a Duration composite.
*/
func fromGo(value time.Duration) Interface {
	return &data{value: value}
}

/*
Parse builds a Duration from a Go duration string (see time.ParseDuration).

On success the Result carries the Duration as its payload; when the input is
malformed the Result carries an Error instead — the operation never panics and
never returns nil.

	r := Duration.Parse("1h30m")
	if r.HasError() {
	    // r.Error().Message()
	}
*/
func Parse(s string) Result.Interface {
	value, err := time.ParseDuration(s)
	if err != nil {
		return Result.New(
			Result.WithError(
				Error.New(err.Error()),
			),
		)
	}
	return Result.New(
		Result.WithPayload(fromGo(value)),
	)
}

/*
ToSeconds returns the duration truncated to a whole number of seconds.
*/
func (d data) ToSeconds() int64 {
	return int64(d.value / time.Second)
}

/*
ToGoString returns the Go textual representation of the duration (e.g. "1h30m0s").
*/
func (d data) ToGoString() string {
	return d.value.String()
}

/*
Add returns a new Duration that is the sum of the receiver and other.
*/
func (d data) Add(other Interface) Interface {
	return fromGo(d.value + time.Duration(other.ToSeconds())*time.Second)
}

/*
Sub returns a new Duration that is the difference of the receiver and other.
*/
func (d data) Sub(other Interface) Interface {
	return fromGo(d.value - time.Duration(other.ToSeconds())*time.Second)
}

/*
Equal reports whether the receiver and other span the same number of seconds.
*/
func (d data) Equal(other Interface) bool {
	return d.ToSeconds() == other.ToSeconds()
}

/*
IsNull reports whether the Duration is the Null-Object variant.

A concrete Duration is never null.
*/
func (d data) IsNull() bool {
	return false
}

type null struct{}

/*
Null returns the Null-Object variant of Duration.

It satisfies Interface so callers never have to test for a bare nil: it spans
zero seconds, its arithmetic returns itself, and IsNull() returns true.
*/
func Null() Interface {
	return &null{}
}

func (n null) ToSeconds() int64 {
	return 0
}

func (n null) ToGoString() string {
	return `<NullDuration>`
}

func (n null) Add(Interface) Interface {
	return &null{}
}

func (n null) Sub(Interface) Interface {
	return &null{}
}

func (n null) Equal(other Interface) bool {
	return other.IsNull()
}

func (n null) IsNull() bool {
	return true
}
