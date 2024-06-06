package parser

import (
	"encoding/json"
	"strings"
	"strconv"
	"fmt"
	"io"
)

func ParseGoTestJson(text string) (*TestSummary, error){
	jsonDecoder := json.NewDecoder(strings.NewReader(text))
	summary := new(TestSummary)
	summary.PackageResults = make(map[string]*PackageResult)

	for {
		var data JsonData = JsonData{}
		if err := jsonDecoder.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			fmt.Errorf("parsing error: %s", err.Error())
			return nil, err
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

	return summary, nil
}