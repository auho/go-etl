package tag

import (
	"os"
	"testing"
)

var _ruleName = "a"
var _ruleTableName = "rule_" + _ruleName
var _rule = &ruleTest{}
var _ruleAliasFixed = &ruleAliasFixedTest{}
var _contents = []string{
	`b一ab一bc一abc一123b一b123一123一0123一1234一01234一`,
	`中文一b中文123一123中文b一中bb文一中123文一中00文一中aa文一中00文一中aa文一中中文文一中二二文一
123一一123一123一123--`,
	`一中文一中b文一中1文一中2文一中3文一中ab文一中12文一中13文一中23文一中123文`,
}

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp()    {}
func tearDown() {}
