package main

import "fmt"
import "flag"

func main() {
    fileNamePtr := flag.String("f", "", "the ownCloud log file")
    flag.Parse()
    fmt.Println("filename:", *fileNamePtr)
}