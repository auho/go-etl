package containment

func ExampleNewContainsAll() {
	NewContainsAll([]string{"1", "2", "12", "ab"}, _rule)
}

func ExampleNewContainsFirstLine() {
	NewContainsFirstLine([]string{"1", "2", "12", "ab"}, _rule)
}
