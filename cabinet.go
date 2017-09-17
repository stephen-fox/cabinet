package cabinet

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// Exists checks if a file exists.
//
// Based on work by "Mostafa": https://stackoverflow.com/a/10510783
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// CopyDirectory copies the source directory into the destination directory.
// If the destination does not exist, the function will automatically
// create it.
func CopyDirectory(sourceDirPath string, destinationDirPath string) error {
	directory, err := os.Stat(sourceDirPath)
	if err != nil {
		return err
	}

	if !directory.IsDir() {
		return errors.New("Failed to copy directory. Source is not a directory")
	}

	err = os.MkdirAll(destinationDirPath, directory.Mode())
	if err != nil {
		return err
	}

	directoryContents, err := ioutil.ReadDir(sourceDirPath)
	if err != nil {
		return err
	}

	for _, content := range directoryContents {
		sourcePath := sourceDirPath + "/" + content.Name()
		if content.IsDir() {
			destinationPath := destinationDirPath + "/" + content.Name()
			err = CopyDirectory(sourcePath, destinationPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(sourcePath, destinationDirPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies the source file into the destination path.
//
// Based on work by "Salvador Dali": https://stackoverflow.com/a/33865286
func CopyFile(sourceFilePath string, destinationDirPath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	fileInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return err
	}

	fullDestinationPath := destinationDirPath + "/" + path.Base(sourceFilePath)
	destFile, err := os.Create(fullDestinationPath)
	if err == nil {
		os.Chmod(fullDestinationPath, fileInfo.Mode())
	} else {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// DownloadFile downloads a file to the destination. The resulting file's name
// must be specified in the download path.
//
// Based on work by "Pablo Jomer": https://stackoverflow.com/a/33845771
func DownloadFile(url string, fileDownloadPath string) error {
	// Create the file.
	out, err := os.Create(fileDownloadPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data.
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// ReplaceLineInFile replaces a line in a file with the desired string. If
// the replacement string already exists, then the function
// immediately returns.
func ReplaceLineInFile(path string, match string, replacement string,
	lineEnding string) (wasReplaced bool, err error) {

	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(contents), lineEnding)
	for i, line := range lines {
		// Skip replacement if the file already contains the replacement.
		if strings.EqualFold(line, replacement) {
			return false, nil
		}
		if strings.Contains(line, match) {
			lines[i] = replacement
		}
	}

	newContents := strings.Join(lines, lineEnding)
	err = ioutil.WriteFile(path, []byte(newContents), fileInfo.Mode())
	if err != nil {
		return false, err
	}

	return true, nil
}
