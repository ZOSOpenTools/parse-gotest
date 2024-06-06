## parse-gotest
A tool to parse results from `go test -v -json`
## Usage 
```
The program parse the output from result go test -json.
                parse-gotest [-options] [files] [directories]

  -d string
        delimeter: set the delimeter for the summary default to new line character (default "\n")
  -detail
        Show the summary and the test output.  Use -type to specify which type of test want to print (default true)
  -list
        list: print only the test name. Use -type to specify which type of test want to list
  -p
        Use with -l or -list. List testcases by package instead of the only print the entire list of testcases
  -summary
        summary: print the total number of test and the number of tests that are pass, fail and skip
  -type string
        Value: all | pass | fail | skip
        Specify which type of test (default "fail")
```
## Install 
build from source code 
```
git clone https://github.com/bluecivet/parse-gotest.git
go build 
```

## Example
Generate a go test result for your program 
```
cd to-your-go-program
go test -v -json ./... > result.txt
```
using the tool
```
parse-gotest result.txt

output:
Test Summary
total: 5, pass: 2, fail: 2, skip: 1


Fail Test:
Fail on:  package: github.com/bluecivet/parse-gotest/test, run: 5, pass: 2, fail: 2, skip: 1
-----------------------------------------------------------

TestFail1
=== RUN   TestFail1
fail1 message
    produce_test.go:18: Fail1 message
--- FAIL: TestFail1 (0.00s)
Elapsed: 0


TestFail2
=== RUN   TestFail2
fail2 message
    produce_test.go:23: Fail2 message
--- FAIL: TestFail2 (0.00s)
Elapsed: 0

-----------------------------------------------------------
```
specify the type of test result you want to see (all, pass, fail, skip)
```
parse-gotest -type skip result.txt

output:
Test Summary
total: 5, pass: 2, fail: 2, skip: 1


Skip Test:
Skip on:  package: github.com/bluecivet/parse-gotest/test, run: 5, pass: 2, fail: 2, skip: 1
-----------------------------------------------------------

TestSkip
=== RUN   TestSkip
    produce_test.go:27: skipping test
--- SKIP: TestSkip (0.00s)
Elapsed: 0

-----------------------------------------------------------
```

print only test summary 
```
parse-gotest -summary  result.txt

output:
total:5
pass:2
fail:2
skip:1
```
you can specify a delimiter for the summary 
```
/parse-gotest -summary -d ", " result.txt

output:
total:5, pass:2, fail:2, skip:1
```

only print a list of test 
```
parse-gotest -list  -type all  result.txt

output:
TestPass1
TestPass2
TestFail1
TestFail2
TestSkip
```

print the list with package information 
```
parse-gotest -list  -type all -p  result.txt

output:
package:  github.com/bluecivet/parse-gotest/test
-----------------------------------------------------------
TestPass1
TestPass2
TestFail1
TestFail2
TestSkip
```
