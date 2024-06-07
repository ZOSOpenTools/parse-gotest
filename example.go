package main

import (
	"fmt"
	"github.com/bluecivet/parse-gotest/parser"
)

func ExampleUsage() {

	testoutput := `{"Time":"2024-05-29T10:50:26.103587893-04:00","Action":"start","Package":"syscall"}
{"Time":"2024-05-29T10:50:26.231295252-04:00","Action":"run","Package":"aaa","Test":"TestExample"}
{"Time":"2024-05-29T10:50:26.231321264-04:00","Action":"output","Package":"aaa","Test":"TestExample","Output":"=== RUN   TestExample\n"}
{"Time":"2024-05-29T10:50:26.231336089-04:00","Action":"output","Package":"aaa","Test":"TestExample","Output":"--- PASS: TestExample (0.00s)\n"}
{"Time":"2024-05-29T10:50:26.231341102-04:00","Action":"pass","Package":"aaa","Test":"TestExample","Elapsed":0}
`
	goTestSummary, err := parser.ParseGoTestJson(testoutput)
	if err != nil {
		println("Cannot parse json data: ", err.Error())
		return
	}

	fmt.Println("Test Summary")
	fmt.Println(goTestSummary)
	fmt.Println("\n")

	fmt.Println("\n========================\n")
	packageList := goTestSummary.TestPackageList()
	fmt.Println("package test: ", packageList)

	for _, pkgName := range packageList {
		fmt.Println("\n----------------------------------\n")
		fmt.Println(goTestSummary.PackageResults[pkgName])
	}

	fmt.Println("\n==================================\n")

	fmt.Println("Fail test:")
	for _, pkgName := range packageList {
		if goTestSummary.PackageResults[pkgName].HasFail() {
			fmt.Println(goTestSummary.PackageResults[pkgName].Summary())
			for _, failTestName := range goTestSummary.PackageResults[pkgName].FailTests {
				fmt.Println(failTestName)
				fmt.Println(goTestSummary.PackageResults[pkgName].TestOutput[failTestName])
			}
		}
	}
}
