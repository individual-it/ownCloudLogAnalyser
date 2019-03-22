package main

import (
    "os"
    "fmt"
    "flag"
    "log"
    "strings"
    "encoding/json"
    "strconv"
    "gopkg.in/Knetic/govaluate.v3"
    "github.com/hpcloud/tail"
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
    if filterResult!=true || err != nil{
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
    fileNamePtr := flag.String("file", "", "the ownCloud log file")
    followPtr := flag.Bool("f", false, "output appended data as the file grows")
    showLineCounterPtr := flag.Bool("linenumbers", false, "show the line numbers")
    listOfKeysToViewStringPtr := flag.String("view", "", "list of keys to be shown (separate by comma), if empty all are shown")
    filterPtr := flag.String("filter", "", "filter the output by logical expressions e.g. \"user=='admin'&&level>=3\"")
    tailPtr := flag.Int64("tail", 0, "show only the n last lines")
    flag.Parse()

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
    var seekOffset int64 = 0
    if *tailPtr > 0 {
        var lineBreaksCounter int64

        logFile, err := os.Open(*fileNamePtr)
        defer logFile.Close()
        checkErr("Failed to read the logfile.", err)
        for lineBreaksCounter = 0; lineBreaksCounter <= *tailPtr; {
            seekOffset--
            _, err = logFile.Seek(seekOffset, os.SEEK_END)
            checkErr("Failed to read the logfile.", err)
            readByte := make([]byte, 1)
            _, err = logFile.Read(readByte)
            checkErr("Failed to read the logfile.", err)
            if string(readByte) == "\n" {
                lineBreaksCounter++
            }
        }
    }
    //main loop to loop through the log file
    logFile, err := tail.TailFile(
        *fileNamePtr,
        tail.Config{
            Follow: *followPtr,
            MustExist: true})
        if *tailPtr > 0 {
            logFile.Location = &tail.SeekInfo{seekOffset+1, os.SEEK_END}
        }

    checkErr("Failed to read the logfile.", err)

    for line := range logFile.Lines {
        currentLogLine := logLine {
            raw: line.Text,
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

    checkErr("Failed to read the logfile.", logFile.Err())
}
