# ownCloud log file analyser

a small tool to analyse the log output of [ownCloud](http://github.com/owncloud/core/)

## Usage:
```     
  -f    output appended data as the file grows
  -file string
        the ownCloud log file
  -filter string
        filter the output by logical expressions e.g. "user=='admin'&&level>=3"
  -linenumbers
        show the line numbers
  -tail int
        show only the n last lines
  -view string
        list of keys to be shown (separate by comma), if empty all are shown      
```

currently the main purpose of this repo is to learn `go`
