package collector_test

import (
	"github.com/auho/go-etl/v3/job/transform/collector"
)

// ExampleNewKeysAll demonstrates creating a collector in semantic A (by key
// position): all keys' contents are searched together.
func ExampleNewKeysAll() {
	// Searches contents of both "name" and "email" together.
	collector.NewKeysAll([]string{"name", "email"}, nil /* extractor */)
}

// ExampleNewKeysFirst demonstrates semantic A: only the first key's content
// is searched, regardless of whether it is empty.
func ExampleNewKeysFirst() {
	// Only "name" (keys[0]) is searched.
	collector.NewKeysFirst([]string{"name", "email"}, nil /* extractor */)
}

// ExampleNewKeysMatchAny demonstrates semantic A: each key is searched one by
// one, stopping at the first match. Empty-value keys are still searched.
func ExampleNewKeysMatchAny() {
	// Searches "phone", then "email", then "wechat" until a match is found.
	collector.NewKeysMatchAny([]string{"phone", "email", "wechat"}, nil /* extractor */)
}

// ExampleNewKeysFirstValued demonstrates semantic B (by value): selects the
// first key with a non-empty value.
// Given keys ["name", "nickname"] and item {"name": "", "nickname": "bob"},
// it searches "bob" (the first non-empty value).
func ExampleNewKeysFirstValued() {
	collector.NewKeysFirstValued([]string{"name", "nickname"}, nil /* extractor */)
}

// ExampleNewKeysMatchAnyValued demonstrates semantic B: iterates only keys
// with non-empty values, searching each one until a match is found.
// Empty-value keys are skipped entirely.
func ExampleNewKeysMatchAnyValued() {
	collector.NewKeysMatchAnyValued([]string{"phone", "email", "wechat"}, nil /* extractor */)
}
