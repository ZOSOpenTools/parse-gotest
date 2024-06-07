package parser

import (
	"fmt"
	"strings"
)

type PackageResult struct {
	Package    string
	RunTests   []string          // names of the tests that are runned. This will be to total tests that are runned
	PassTests  []string          // names of the tests that are passed
	FailTests  []string          // names of the tests that are failed
	SkipTests  []string          // names of the tests that are skiped
	TestOutput map[string]string // output of each test results: test name -> testcase output
}

func NewPackageResult(pkgName string) *PackageResult {
	pr := new(PackageResult)
	pr.Package = pkgName
	pr.RunTests = make([]string, 0, 10)
	pr.PassTests = make([]string, 0, 10)
	pr.FailTests = make([]string, 0, 10)
	pr.SkipTests = make([]string, 0, 10)
	pr.TestOutput = make(map[string]string)
	return pr
}

func (pr *PackageResult) IsAllPass() bool {
	return len(pr.RunTests) == len(pr.PassTests)
}

func (pr *PackageResult) HasFail() bool {
	return len(pr.FailTests) > 0
}

func (pr *PackageResult) HasSkip() bool {
	return len(pr.SkipTests) > 0
}

func (pr *PackageResult) Summary() (output string) {
	run := len(pr.RunTests)
	pass := len(pr.PassTests)
	fail := len(pr.FailTests)
	skip := len(pr.SkipTests)
	output = fmt.Sprintf("package: %s, run: %d, pass: %d, fail: %d, skip: %d",
		pr.Package,
		run,
		pass,
		fail,
		skip)
	return
}

func (pr *PackageResult) AllTestStatus() map[string]string {
	result := make(map[string]string)
	// init all test to unknow status
	// e.g: setting "testA": "?"
	for _, runTest := range pr.RunTests {
		result[runTest] = "?"
	}

	// e.g: setting "testA": "pass | fail | skip"
	setTestStatus(&result, &pr.PassTests, "pass")
	setTestStatus(&result, &pr.FailTests, "fail")
	setTestStatus(&result, &pr.SkipTests, "skip")
	return result
}

func setTestStatus(testSet *map[string]string, statusTestSet *[]string, status string) {
	// testSet must be init first
	for _, test := range *statusTestSet {
		if _, ok := (*testSet)[test]; ok {
			(*testSet)[test] = status
		} else {
			(*testSet)[test] = "? (test is not in runned test set)"
		}
	}
}

func (pr *PackageResult) String() string {
	summaryStr := pr.Summary()
	var sb strings.Builder
	sb.WriteString(summaryStr)

	allTestStatus := pr.AllTestStatus()
	for testName, testStatus := range allTestStatus {
		// output:
		// \t testName: ok | fail | skip | ?
		// "?" for test in unknow result
		sb.WriteString("\t")
		sb.WriteString(testName)
		sb.WriteString(": ")
		sb.WriteString(testStatus)
		sb.WriteString("\n")
	}
	return sb.String()

}
