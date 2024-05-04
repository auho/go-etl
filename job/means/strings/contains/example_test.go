package contains

func ExampleNewContainsAll() {
	NewContainsAll(_rule, []string{"1", "2", "12", "ab"}, NewExportAll(_rule))
}

func ExampleNewContainsFirst() {
	NewContainsFirst(_rule, []string{"1", "2", "12", "ab"}, NewExportLine(_rule))
}
