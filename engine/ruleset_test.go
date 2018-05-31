package engine

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
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

func TestExpandMultiLevelImport(t *testing.T) {
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

func TestExpandMultiImport(t *testing.T) {
	// setup
	import2FilePath, err := writeTempFile([]byte(`you are looking for...`))
	assert.NoError(t, err, "test file write error")
	defer os.Remove(import2FilePath) // clean up

	import1FilePath, err := writeTempFile([]byte(`the imports`))
	assert.NoError(t, err, "test file write error")
	defer os.Remove(import1FilePath) // clean up

	rules := []byte(fmt.Sprintf(`These are not
		{@ %s @}
		{@ %s @}`, import1FilePath, import2FilePath))
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
	rules := []byte(`no imports to see here`)

	ruleset := new(RuleSet)
	ruleset.SetImportDepth(3)
	expandedRules, err := ruleset.expandImports(rules, 0)

	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, string(rules), string(expandedRules), "rules do not match")
}
func TestBadImport(t *testing.T) {
	rules := []byte(`No import {@ nomanland @}`)

	ruleset := new(RuleSet)
	ruleset.SetImportDepth(3)
	expandedRules, err := ruleset.expandImports(rules, 0)

	assert.NoError(t, err, "unexpected error")

	assert.Equal(t, "No import ", string(expandedRules), "rules do not match")
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

func TestRelativePath(t *testing.T) {
	importFilePath, err := writeTempFile([]byte(`android two`))
	assert.NoError(t, err, "test file write error")
	defer os.Remove(importFilePath) // clean up

	importFileName := path.Base(importFilePath)

	rules := []byte(fmt.Sprintf(`android one
		{@ %s @}`, importFileName))
	rulesFilePath, err := writeTempFile(rules)
	assert.NoError(t, err, "test file write error")
	defer os.Remove(rulesFilePath) // clean up

	expectedRules := []byte(`android one
		android two`)

	ruleset, err := NewRuleset(rulesFilePath)
	ruleset.SetImportDepth(3)

	expandedRules, err := ruleset.expandImports(rules, 0)

	assert.NoError(t, err, "unexpected error")

	assert.NotEqual(t, string(rules), string(expandedRules), "no changes made")
	assert.Equal(t, string(expectedRules), string(expandedRules), "rules do not match")
}

func TestDirectoryList(t *testing.T) {
	importDirPath, err := ioutil.TempDir("", "templr-test")
	assert.NoError(t, err, "failed to make dir")
	defer os.RemoveAll(importDirPath) // clean up

	expectedFileList := []string{}
	fileCount := 3
	for i := 0; i < fileCount; i++ {
		fileInfo, err := ioutil.TempFile(importDirPath, "file")
		assert.NoError(t, err, "")
		if err != nil {
			continue

		}
		expectedFileList = append(expectedFileList, fileInfo.Name())
	}
	symlinkPath := path.Join(importDirPath, "link")
	os.Symlink(expectedFileList[0], symlinkPath)
	expectedFileList = append(expectedFileList, symlinkPath)
	for i := 0; i < 3; i++ {
		tmpDirpath, err := ioutil.TempDir(importDirPath, "dir")
		assert.NoError(t, err, "error making dir")

		_, err = ioutil.TempFile(tmpDirpath, "file")
		assert.NoError(t, err, "error making file")
	}

	ruleset := new(RuleSet)
	ruleset.SetImportDepth(3)
	resultFileList, err := ruleset.getFileList(path.Join(importDirPath, "*"))
	assert.NoError(t, err, "unknown")

	assert.Equal(t, len(expectedFileList), len(resultFileList))
	for i := range expectedFileList {
		assert.Contains(t, resultFileList, expectedFileList[i], "expected item not found")
	}
}

func writeTempFile(contents []byte) (name string, err error) {
	fileObj, err := ioutil.TempFile(os.TempDir(), "templr-test")
	if err != nil {
		return "", errors.Wrapf(err, "could not open file to write")
	}
	defer fileObj.Close()

	if _, err = io.Copy(fileObj, bytes.NewReader(contents)); err != nil {
		return "", errors.Wrapf(err, "could not write to file")
	}
	return fileObj.Name(), nil
}
