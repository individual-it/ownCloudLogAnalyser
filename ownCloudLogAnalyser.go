package main

import (
    "fmt"
    "flag"
    "log"
    "os"
    "bufio"
)

func main() {
    fileNamePtr := flag.String("f", "", "the ownCloud log file")
    showLineCounterPtr := flag.Bool("linenumbers", false, "show the line numbers")
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
    for logFileScanner.Scan() {
        if *showLineCounterPtr {
            fmt.Printf("%v: %v\n", lineCounter, logFileScanner.Text())
        } else {
            fmt.Println(logFileScanner.Text())
        }
        lineCounter++
    }

    if err := logFileScanner.Err(); err != nil {
        log.Fatal(err)
    }
}