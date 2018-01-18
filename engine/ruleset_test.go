package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractVarsWithValidVars(t *testing.T) {
	ruleset := new(RuleSet)

	rules := `this
	{$ test: ['one', 'two', 'three'] $}
	is
	{$ test2: ['one', 'two', 'three']$}normal
	{$ test3: 
   - 'one'
   - 'two'
   - 'three'
	$}
	text`
	ruleBytes := []byte(rules)

	expectedRules := []byte(`this
	isnormal
	text`)
	varTmp := []interface{}{"one", "two", "three"}
	expectedVars := map[string]interface{}{"test": varTmp, "test2": varTmp, "test3": varTmp}

	resultRules, resultVars, err := ruleset.extractVars(ruleBytes)
	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, expectedVars, resultVars, "vars do not match")
	assert.Equal(t, expectedRules, resultRules, "rules do not match")

}

func TestExtractVarsWithNoVars(t *testing.T) {
	rules := []byte(`this is
void of rules`)

	ruleset := new(RuleSet)
	resultRules, resultVars, err := ruleset.extractVars(rules)

	assert.NoError(t, err, "unexpected error")

	assert.Empty(t, resultVars, "unexpected vars value")
	assert.Equal(t, rules, resultRules, "rules do not match")

}
