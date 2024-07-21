package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func ParseGoTestJson(text string) (*TestSummary, error) {
	scanner := bufio.NewScanner(strings.NewReader(text))
	summary := new(TestSummary)
	summary.PackageResults = make(map[string]*PackageResult)
	errs := make([]error, 0, 3)
	for scanner.Scan() {
		line := scanner.Text()
		var data JsonData = JsonData{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			err = fmt.Errorf("parsing error: %s, line: %s", err.Error(), line)
			errs = append(errs, err)
		}

		pkgName := data.Package
		if pkgName == "" {
			println("package name is not string")
			continue
		}

		if _, ok := summary.PackageResults[pkgName]; !ok {
			summary.PackageResults[pkgName] = NewPackageResult(pkgName)
		}

		if test := data.Test; test != "" {
			// collect the test output data
			if output := data.Output; output != "" {
				summary.PackageResults[pkgName].TestOutput[test] += output
			} else if elapsed := data.Elapsed; elapsed >= 0 {
				summary.PackageResults[pkgName].TestOutput[test] += "Elapsed: " + strconv.FormatFloat(elapsed, 'g', -1, 64) + "\n"
			} else {
				if _, ok := summary.PackageResults[pkgName].TestOutput[test]; ok {
					summary.PackageResults[pkgName].TestOutput[test] += test + "\n"
				} else {
					summary.PackageResults[pkgName].TestOutput[test] = test + "\n"
				}
			}

			// depend on the action record the test
			if action := data.Action; action != "" {
				switch action {
				case "output", "pause", "cont":
					// do nothing. output is handle by above
				case "run":
					summary.Total++
					summary.PackageResults[pkgName].RunTests = append(summary.PackageResults[pkgName].RunTests, test)
				case "pass":
					summary.Pass++
					summary.PackageResults[pkgName].PassTests = append(summary.PackageResults[pkgName].PassTests, test)
				case "fail":
					summary.Fail++
					summary.PackageResults[pkgName].FailTests = append(summary.PackageResults[pkgName].FailTests, test)
				case "skip":
					summary.Skip++
					summary.PackageResults[pkgName].SkipTests = append(summary.PackageResults[pkgName].SkipTests, test)
				default:
					println("unknow action: ", action)
				}
			}

		}
	}

	// unknow action view as fail
	// in go1.22 running "go test -v -json" for a panic
	// it causes the package to fail but not show test action as fail
	for _, pr := range summary.PackageResults {
		summary.Fail += handleExtraTest(pr)
	}

	if len(errs) > 0 {
		var sb strings.Builder
		for _, err := range errs {
			sb.WriteString(err.Error())
			sb.WriteString("\n")
		}
		return summary, fmt.Errorf("parse json error: \n%s", sb.String())
	}

	return summary, nil
}

// If a test is run but does not produce a result
// (pass, fail, skip) in JSON file, we mark it as a failed test.
// An example of this is a test that causes a panic,
// which can happen in Go 1.22.
func handleExtraTest(pr *PackageResult) (numOfFail int) {
	if len(pr.RunTests) == len(pr.PassTests)+len(pr.FailTests)+len(pr.SkipTests) {
		return
	}

	visitMap := make(map[string]bool, len(pr.RunTests))
	// init visitMap
	for _, runTest := range pr.RunTests {
		visitMap[runTest] = false
	}

	// start visiting
	for _, testName := range pr.PassTests {
		visitMap[testName] = true
	}

	for _, testName := range pr.FailTests {
		visitMap[testName] = true
	}

	for _, testName := range pr.SkipTests {
		visitMap[testName] = true
	}

	// check result
	for testName, visited := range visitMap {
		if !visited {
			numOfFail++
			pr.FailTests = append(pr.FailTests, testName)
		}
	}
	return
}
