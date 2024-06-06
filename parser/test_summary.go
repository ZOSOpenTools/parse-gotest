package parser

import (
	"fmt"
)

type TestSummary struct{
	Pass int
	Fail int
	Skip int
	Total int

	// each package result: package name -> the result of the package
	PackageResults map[string]*PackageResult
}

func (ts *TestSummary) IsAllPass() bool {
	return ts.Pass == ts.Total
}

func (ts *TestSummary) HasFail() bool {
	return ts.Fail > 0
}

// return package names that the test runned
func (ts *TestSummary) TestPackageList() []string {
	result := make([]string, 0, 3)
	for pkgName, _ := range ts.PackageResults {
		result = append(result, pkgName)
	}
	return result
}

func (ts *TestSummary) String() string {
	return fmt.Sprintf("total: %d, pass: %d, fail: %d, skip: %d", ts.Total, ts.Pass, ts.Fail, ts.Skip)
}