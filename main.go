package main

import (
	"fmt"
	"time"

	Result "github.com/go-composites/result/src"
	Time "github.com/go-composites/time/src"
	Duration "github.com/go-composites/time/src/duration"
)

func main() {
	epoch := Time.FromUnix(0)
	later := Time.FromUnix(3600)

	fmt.Println("epoch:", epoch.ToGoString())
	fmt.Println("later:", later.ToGoString())

	// Comparisons are plain Go bools.
	fmt.Println("epoch before later:", epoch.Before(later))

	// Sub yields the Duration between two Times.
	gap := later.Sub(epoch)
	fmt.Println("gap:", gap.ToGoString()) // 1h0m0s

	// Add returns a Result — shifting an instant by a Duration.
	shifted := epoch.Add(Duration.FromSeconds(90))
	fmt.Println("epoch + 90s:", shifted.Payload().(Time.Interface).ToGoString())

	// Parsing is fallible: a malformed value is a value, not a panic.
	bad := Time.Parse(time.RFC3339, "not-a-time")
	fmt.Println("bad parse has error:", bad.HasError())
	fmt.Println(bad.Error().Message())

	good := Duration.Parse("1h30m")
	report("1h30m", good)

	// The Null-Object variant is never a bare nil.
	fmt.Println("null time is null:", Time.Null().IsNull())
}

func report(label string, result Result.Interface) {
	if result.HasError() {
		fmt.Printf("%s -> error: %s\n", label, result.Error().Message())
		return
	}
	fmt.Printf("%s -> %d seconds\n", label, result.Payload().(Duration.Interface).ToSeconds())
}
