package regexps

func ExampleNewAllSubMatch() {
	NewAllSubMatch([]string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, _rule)
}

func ExampleNewSubMatchAllLine() {
	NewSubMatchAllLine([]string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, _rule)
}

func ExampleNewSubMatchFirstFlag() {
	NewSubMatchFirstFlag([]string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, _rule)
}
