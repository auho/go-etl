package contains

func ExampleNewContainsAll() {
	NewContainsAll([]string{"1", "2", "12", "ab"}, NewExportAll(_rule))
}

func ExampleNewContainsFirst() {
	NewContainsFirst([]string{"1", "2", "12", "ab"}, NewExportLine(_rule))
}
