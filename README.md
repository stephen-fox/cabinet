# cabinet

## What is it?
A Go library that provides several helpful functions for working with files.

## How do I use it?
Execute the following on the commandline:
```
go get github.com/stephen-fox/cabinet
```

Then, in your Go application, add the following import statement:
```
import (
    "github.com/stephen-fox/cabinet"
)
```

## What does the API look like?

* `Exists` - Check if a file or directory exists
* `FileExists` - Check if a file exists
* `DirectoryExists` - Check if a directory exists
* `CopyFilesWithSuffix` - Recursively copy files ending with a suffix.
Optionally specify if existing files should be overwritten
* `CopyDirectory` - Recursively copy a directory. Optionally specify if
existing files should be overwritten
* `CopyFile` - Copy a file. Optionally specify if an existing file should
be overwritten
* `DownloadFile` - Download a file
* `ReplaceLineInFile` - Replace a line in a file
* `GetFileHash` - Get a file's hash using the `hash.Hash` interface