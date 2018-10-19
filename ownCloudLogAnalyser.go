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

func main() {
    fileNamePtr := flag.String("f", "", "the ownCloud log file")
    showLineCounterPtr := flag.Bool("linenumbers", false, "show the line numbers")
    listOfKeysToViewStringPtr := flag.String("view", "", "list of keys to be shown (separate by comma), if empty all are shown")
    filterPtr := flag.String("filter", "", "filter the output by logical expressions e.g. \"user=='admin'&&level>=3\"")
    flag.Parse()

    logFile, err := os.Open(*fileNamePtr)
    if err != nil {
        log.Fatal(err)
    }
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
        //parse unstructured json
        //see https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/
        var decodedData map[string]interface{}
        if needJsonDecode {
            err := json.Unmarshal([]byte(logFileScanner.Text()), &decodedData)
            if err != nil {
                log.Fatal("JSON error, line: ", strconv.Itoa(lineCounter), " error: ", err)
            }
        }
        filterExpression, err := govaluate.NewEvaluableExpression(*filterPtr);
        if err != nil {
            log.Fatal("cannot evaluate filter string: ", err)
        }
        filterResult, err := filterExpression.Evaluate(decodedData);
        if err != nil {
            log.Fatal("cannot evaluate filter string: ", err)
        }

        if filterResult==true {
            if *showLineCounterPtr {
                fmt.Printf("%v: ", lineCounter)
            }
            if len(listOfKeysToView) > 0 {
                //show only the keys we requested
                for _, jsonKey := range listOfKeysToView {
                    fmt.Printf("%v: %v\t", jsonKey, decodedData[jsonKey])
                }
                fmt.Println()
            } else {
                fmt.Println(logFileScanner.Text())
            }
        }
        lineCounter++
    }

    if err := logFileScanner.Err(); err != nil {
        log.Fatal(err)
    }
}