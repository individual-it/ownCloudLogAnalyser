package main

import (
    "fmt"
    "flag"
    "log"
    "os"
    "bufio"
    "strings"
    "encoding/json"
    "strconv"
    "gopkg.in/Knetic/govaluate.v3"
)

type logLine struct {
    lineNumber int
    raw string
    decodedData map[string]interface{}
    showLine bool
}

// decodes the raw data and updates logLine.decodedData
func (logLine *logLine) jsonDecode() {
    //parse unstructured json
    //see https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/
    err := json.Unmarshal([]byte(logLine.raw), &logLine.decodedData)
    checkErr("JSON error in line: " + strconv.Itoa(logLine.lineNumber) + ".", err)
}

// evaluates the given filter and updates logLine.showLine according to the filter 
func (logLine *logLine) evaluateFilter(filter string) {
    filterExpression, err := govaluate.NewEvaluableExpression(filter);
    checkErr("Cannot evaluate filter string.", err)
    filterResult, err := filterExpression.Evaluate(logLine.decodedData);
    checkErr("Cannot evaluate filter string.", err)
    if filterResult!=true {
        logLine.showLine = false
    }
}

// prints the formated output
func (logLine *logLine) printLine(showLineCounter bool, listOfKeysToView []string) {
    if logLine.showLine {
        if showLineCounter {
            fmt.Printf("%v: ", logLine.lineNumber)
        }
        if len(listOfKeysToView) > 0 {
            //show only the keys we requested
            for _, jsonKey := range listOfKeysToView {
                fmt.Printf("%v: %v\t", jsonKey, logLine.decodedData[jsonKey])
            }
            fmt.Println()
        } else {
            fmt.Println(logLine.raw)
        }
    }
}

func checkErr(message string, err error) {
    //for error handling see https://davidnix.io/post/error-handling-in-go/
    if err != nil {
        log.Fatal("ERROR! ", message, " Details: ", err)
    }
}

func main() {
    fileNamePtr := flag.String("f", "", "the ownCloud log file")
    showLineCounterPtr := flag.Bool("linenumbers", false, "show the line numbers")
    listOfKeysToViewStringPtr := flag.String("view", "", "list of keys to be shown (separate by comma), if empty all are shown")
    filterPtr := flag.String("filter", "", "filter the output by logical expressions e.g. \"user=='admin'&&level>=3\"")
    flag.Parse()

    logFile, err := os.Open(*fileNamePtr)
    checkErr("Failed to read the logfile.", err)

    //make sure the file is closed also when somethig fails
    //see https://blog.golang.org/defer-panic-and-recover
    defer logFile.Close()

    logFileScanner := bufio.NewScanner(logFile)
    lineCounter := 1
    needJsonDecode := false

    //split the list of keys to view into slices
    listOfKeysToView := make([]string,0)
    if *listOfKeysToViewStringPtr != "" || *filterPtr != "" {
       needJsonDecode = true
    }

    if *listOfKeysToViewStringPtr != "" {
       listOfKeysToView = strings.Split(*listOfKeysToViewStringPtr, ",")
    }

    //main loop to loop through the log file
    for logFileScanner.Scan() {
        currentLogLine := logLine {
            raw: logFileScanner.Text(),
            lineNumber: lineCounter,
            showLine: true }
        if needJsonDecode {
            currentLogLine.jsonDecode()
        }
        if *filterPtr != "" {
            currentLogLine.evaluateFilter(*filterPtr)
        }

        currentLogLine.printLine(*showLineCounterPtr, listOfKeysToView)
        lineCounter++
    }

    checkErr("Failed to read the logfile.", logFileScanner.Err())
}