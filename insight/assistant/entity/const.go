package entity

// Table name prefixes
const (
	NameData       = "data"        // prefix for processed data tables
	NameTag        = "tag"         // prefix for tag/tagging result tables
	NameRule       = "rule"        // prefix for rule tables
	NameDeleted    = "deleted"     // prefix for deleted-row tables
	NameSegWords   = "seg_words"   // suffix for segmented words tables
	NameSplitWords = "split_words" // suffix for split words tables
)

// Column name suffixes for rule-related tables
const (
	NameLabelNum      = "label_num"      // label count column
	NameKeyword       = "keyword"        // keyword column
	NameKeywordNum    = "keyword_num"    // keyword count column
	NameKeywordAmount = "keyword_amount" // keyword total amount column
	NameKeywordLen    = "keyword_len"    // keyword length column
)

// Common column names
const (
	NameWord = "word" // word content column
	NameFlag = "flag" // flag column
	NameNum  = "num"  // numeric count column
)
