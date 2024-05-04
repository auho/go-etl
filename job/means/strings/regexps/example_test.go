package regexps

func ExampleNewAllSubMatch() {
	NewAllSubMatch([]string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, NewExportAll(_rule))
}

func ExampleNewSubMatchAll() {
	NewSubMatchAll([]string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, NewExportLine(_rule))
}

func ExampleNewSubMatchFirst() {
	NewSubMatchFirst([]string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, NewExportFlag(_rule))
}
