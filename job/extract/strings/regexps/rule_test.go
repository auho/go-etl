package regexps

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/extract/testutil"
)

type ruleTest = testutil.RuleTest
type ruleAliasFixedTest = testutil.RuleAliasFixedTest

var _ extract.Rule = (*ruleTest)(nil)
var _ extract.Rule = (*ruleAliasFixedTest)(nil)