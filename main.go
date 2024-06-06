package main

import (
	"fmt"
	"strings"
	"path/filepath"
	"os"
	"flag"

	"github.com/bluecivet/parse-gotest/parser"
)

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	return fileInfo.IsDir() 
}

func readFile(path string) string {
	text, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(text)
}

func readInput(args []string) string {
	currentPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var sb strings.Builder
	for _, arg := range args {
		targetPath := filepath.Join(currentPath, arg)

		if isDir(targetPath) {
			files, err := os.ReadDir(targetPath)
			if err != nil {
				panic(err)
			}

			for _, fileEntry := range files {
				filePath := filepath.Join(targetPath, fileEntry.Name())
				if fileEntry.IsDir() {
					continue
				}
				sb.WriteString(readFile(filePath))
			}
		} else { // file is not a directory
			sb.WriteString(readFile(targetPath))
		}
	}
	return sb.String()
}

func outputFailTestDetail(goTestSummary *parser.TestSummary) {
	fmt.Println("Test Summary")
	fmt.Println(goTestSummary)
	fmt.Println("\n")

	packageList := goTestSummary.TestPackageList()

	fmt.Println("Fail Test:")
	for _, pkgName := range packageList {
		if goTestSummary.PackageResults[pkgName].HasFail() {
			fmt.Println("Fail on: ", goTestSummary.PackageResults[pkgName].Summary())
			fmt.Println("-----------------------------------------------------------")
			for _, failTestName := range goTestSummary.PackageResults[pkgName].FailTests {
				fmt.Println()
				fmt.Println(goTestSummary.PackageResults[pkgName].TestOutput[failTestName])
			}
			fmt.Println("-----------------------------------------------------------")
		}
	}
}

func outputAllTestDetail(goTestSummary *parser.TestSummary) {
	fmt.Println("Test Summary")
	fmt.Println(goTestSummary)
	fmt.Println("\n")

	packageList := goTestSummary.TestPackageList()

	fmt.Println("Test Package:", packageList)
	for _, pkgName := range packageList {
		fmt.Println("-----------------------------------------------------------")
		fmt.Println( goTestSummary.PackageResults[pkgName])
		fmt.Println("-----------------------------------------------------------")
	}
}

func outputPassTestDetail(goTestSummary *parser.TestSummary) {
	fmt.Println("Test Summary")
	fmt.Println(goTestSummary)
	fmt.Println("\n")

	packageList := goTestSummary.TestPackageList()

	fmt.Println("Skip Test:")
	for _, pkgName := range packageList {
		if len(goTestSummary.PackageResults[pkgName].PassTests) > 0 {
			fmt.Println(goTestSummary.PackageResults[pkgName].Summary())
			fmt.Println("-----------------------------------------------------------")
			for _, passTestName := range goTestSummary.PackageResults[pkgName].PassTests {
				fmt.Println()
				fmt.Println(goTestSummary.PackageResults[pkgName].TestOutput[passTestName])
			}
			fmt.Println("-----------------------------------------------------------")
		}
	}
}

func outputSkipTestDetail(goTestSummary *parser.TestSummary) {
	fmt.Println("Test Summary")
	fmt.Println(goTestSummary)
	fmt.Println("\n")

	packageList := goTestSummary.TestPackageList()

	fmt.Println("Skip Test:")
	for _, pkgName := range packageList {
		if goTestSummary.PackageResults[pkgName].HasSkip() {
			fmt.Println("Skip on: ", goTestSummary.PackageResults[pkgName].Summary())
			fmt.Println("-----------------------------------------------------------")
			for _, skipTestName := range goTestSummary.PackageResults[pkgName].SkipTests {
				fmt.Println()
				fmt.Println(goTestSummary.PackageResults[pkgName].TestOutput[skipTestName])
			}
			fmt.Println("-----------------------------------------------------------")
		}
	}
}

func outputTestName(goTestSummary *parser.TestSummary, testType string, printByPackage bool) {
	var tests []string
	for pkgName, pr := range goTestSummary.PackageResults {
		switch testType {
			case "all":  tests =append(tests, pr.RunTests...)
			case "pass": tests =append(tests, pr.PassTests...)
			case "fail": tests =append(tests, pr.FailTests...)
			case "skip": tests =append(tests, pr.SkipTests...)
			default:
				println("unknow test type: ", testType)
				return
		}
		if printByPackage && len(tests) > 0{
			fmt.Println("package: ", pkgName)
			fmt.Println("-----------------------------------------------------------")
			fmt.Println(strings.Join(tests, "\n"), "\n")
			tests = tests[:0]   // clear tests for next package
		}
	}
	
	if !printByPackage {
		// tests contains all the test name from all the package in here
		fmt.Println(strings.Join(tests, "\n"))
	}
}

func outputSummary(goTestSummary *parser.TestSummary, delimeter string) {
	fmt.Printf("total:%d%spass:%d%sfail:%d%sskip:%d\n", 
	goTestSummary.Total,
	delimeter,
	goTestSummary.Pass,
	delimeter,
	goTestSummary.Fail,
	delimeter,
	goTestSummary.Skip)
}

var (
	summary *bool 
	delimeter  *string
	testType *string 
	detail *bool 
	list *bool
	listByPackage *bool
)

func init() {
	flag.Usage = func() {
		helpText := `The program parse the output from result go test -json.
		parse-gotest [-options] [files] [directories]` 
		fmt.Println(helpText, "\n")
		flag.PrintDefaults()
	}
	summary = flag.Bool("summary", false, "summary: print the total number of test and the number of tests that are pass, fail and skip")
	delimeter = flag.String("d", "\n", "delimeter: set the delimeter for the summary default to new line character")
	testType = flag.String("type", "fail", "Value: all | pass | fail | skip\nSpecify which type of test")
	detail = flag.Bool("detail", true, "Show the summary and the test output.  Use -type to specify which type of test want to print")
	list = flag.Bool("list", false, "list: print only the test name. Use -type to specify which type of test want to list")
	listByPackage = flag.Bool("p", false, "\nUse with -l or -list. List testcases by package instead of the only print the entire list of testcases")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("please pass a file or directory as arguments")
		return
	}

	// handle differently for -detail because by default it is enable 
	// if -list or -summary is specify, -detail need to be specify explicitly to be enable
	isSetDetail := false
	isSetOther := false
	flag.Visit(func (f *flag.Flag) {
		switch f.Name {
			case "summary", "s", "list", "l": isSetOther = true 
			case "detail": isSetDetail = true
		}
	})
	if isSetOther && !isSetDetail {
		*detail = false
	}

	testoutput := readInput(args)

	goTestSummary, err := parser.ParseGoTestJson(testoutput)
	if err != nil {
		println("Cannot parse json data: ", err.Error())
		return
	}

	if *summary {
		outputSummary(goTestSummary, *delimeter)
	}

	if *detail {
		switch *testType {
			case "all": outputAllTestDetail(goTestSummary)
			case "fail": outputFailTestDetail(goTestSummary)
			case "skip": outputSkipTestDetail(goTestSummary)
			case "pass": outputPassTestDetail(goTestSummary)
			default: println("unknow tyoe: ", testType)
		}
	}

	if *list {
		outputTestName(goTestSummary, *testType, *listByPackage)
	}
	
}

