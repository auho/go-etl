package regexps

func ExampleNewAllSubMatch() {
	NewAllSubMatch(_rule, []string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, NewExportAll(_rule))
}

func ExampleNewSubMatchAll() {
	NewSubMatchAll(_rule, []string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, NewExportLine(_rule))
}

func ExampleNewSubMatchFirst() {
	NewSubMatchFirst(_rule, []string{
		`a.{1,2}c`,
		`\b(1)\b`,
		`\b(a)\b`,
		`.*(ab).*`,
	}, NewExportFlag(_rule))
}
