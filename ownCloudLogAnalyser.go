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
    flag.Parse()

    logFile, err := os.Open(*fileNamePtr)
    if err != nil {
        log.Fatal(err)
    }
    //make sure the file is closed also when somethig fails
    //see https://blog.golang.org/defer-panic-and-recover
    defer logFile.Close()

    logFileScanner := bufio.NewScanner(logFile)
    for logFileScanner.Scan() {
        fmt.Println(logFileScanner.Text())
    }

    if err := logFileScanner.Err(); err != nil {
        log.Fatal(err)
    }
}