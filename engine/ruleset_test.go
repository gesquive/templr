package engine

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pkg/errors"
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

func TestExpandSingleImport(t *testing.T) {
	// setup
	importFilePath, err := writeTempFile([]byte(`<contents of test_import>`))
	assert.NoError(t, err, "test file write error")
	defer os.Remove(importFilePath) // clean up

	rules := []byte(fmt.Sprintf(`before import
	{@ %s @}
	after import`, importFilePath))
	expectedRules := []byte(`before import
	<contents of test_import>
	after import`)

	ruleset := new(RuleSet)
	ruleset.SetImportDepth(3)
	expandedRules, err := ruleset.expandImports(rules, 0)

	assert.NoError(t, err, "unexpected error")

	assert.NotEqual(t, string(rules), string(expandedRules), "no changes made")
	assert.Equal(t, string(expectedRules), string(expandedRules), "rules do not match")
}

func TestExpandMultiImport(t *testing.T) {
	// setup
	importLvl3FilePath, err := writeTempFile([]byte(`you are looking for...`))
	assert.NoError(t, err, "test file write error")
	defer os.Remove(importLvl3FilePath) // clean up

	importLvl2FilePath, err := writeTempFile([]byte(fmt.Sprintf(`the imports
		{@ %s @}`, importLvl3FilePath)))
	assert.NoError(t, err, "test file write error")
	defer os.Remove(importLvl2FilePath) // clean up

	rules := []byte(fmt.Sprintf(`These are not
		{@ %s @}`, importLvl2FilePath))
	expectedRules := []byte(`These are not
		the imports
		you are looking for...`)

	ruleset := new(RuleSet)
	ruleset.SetImportDepth(3)
	expandedRules, err := ruleset.expandImports(rules, 0)

	assert.NoError(t, err, "unexpected error")

	assert.NotEqual(t, string(rules), string(expandedRules), "no changes made")
	assert.Equal(t, string(expectedRules), string(expandedRules), "rules do not match")
}

func TestNoImport(t *testing.T) {
	// setup
	importFilePath, err := writeTempFile([]byte(`<contents of test_import>`))
	assert.NoError(t, err, "test file write error")
	defer os.Remove(importFilePath) // clean up

	rules := []byte(`no imports to see here`)

	ruleset := new(RuleSet)
	ruleset.SetImportDepth(3)
	expandedRules, err := ruleset.expandImports(rules, 0)

	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, string(rules), string(expandedRules), "rules do not match")
}

func TestMaxDepthImports(t *testing.T) {
	// setup
	importFilePath, err := writeTempFile([]byte(`{@ file.txt @}`))
	assert.NoError(t, err, "test file write error")
	defer os.Remove(importFilePath) // clean up

	rules := []byte(fmt.Sprintf(`before import
	{@ %s @}
	after import`, importFilePath))
	expectedRules := []byte(`before import
	
	after import`)

	ruleset := new(RuleSet)
	ruleset.SetImportDepth(1)
	expandedRules, err := ruleset.expandImports(rules, 0)

	assert.NoError(t, err, "unexpected error")

	assert.NotEqual(t, string(rules), string(expandedRules), "no changes made")
	assert.Equal(t, string(expectedRules), string(expandedRules), "rules do not match")
}

func writeTempFile(contents []byte) (name string, err error) {
	fileObj, err := ioutil.TempFile(os.TempDir(), "shield-test")
	if err != nil {
		return "", errors.Wrapf(err, "could not open file to write")
	}
	defer fileObj.Close()

	if _, err = io.Copy(fileObj, bytes.NewReader(contents)); err != nil {
		return "", errors.Wrapf(err, "could not write to file")
	}
	return fileObj.Name(), nil
}
