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

// Exists checks if a file or directory exists.
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

// FileExists checks if a file exists.
//
// Based on work by "Mostafa": https://stackoverflow.com/a/10510783
func FileExists(path string) bool {
	file, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	if file.IsDir() {
		return false
	}
	return true
}

// DirectoryExists checks if a directory exists.
//
// Based on work by "Mostafa": https://stackoverflow.com/a/10510783
func DirectoryExists(path string) bool {
	file, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	if !file.IsDir() {
		return false
	}
	return true
}

// CopyFilesWithSuffix recursively copies all files ending with a specific
// suffix into the destination path. If the destination does not exist, the
// function will automatically create it. For example, given a source of
// '/home/user1' and a destination of '/tmp/hello', the resulting directory
// would be '/tmp/hello'.
func CopyFilesWithSuffix(sourceDirPath string, destinationDirPath string, suffix string, shouldOverwrite bool) error {
	contents, err := ioutil.ReadDir(sourceDirPath)
	if err != nil {
		return err
	}

	for _, content := range contents {
		filePath := sourceDirPath + "/" + content.Name()
		if content.IsDir() {
			nextDir := destinationDirPath + "/" + content.Name()
			err := CopyFilesWithSuffix(filePath, nextDir, suffix, shouldOverwrite)
			if err != nil {
				return err
			}
		} else {
			if strings.HasSuffix(content.Name(), suffix) {
				err = CopyFile(filePath, destinationDirPath, shouldOverwrite)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// CopyDirectory recursively copies the source directory's contents into the
// destination directory. If the destination does not exist, the function will
// automatically create it. For example, given a source directory of
// '/home/user1' and a destination of '/tmp/junk', the resulting directory
// would be '/tmp/junk'.
func CopyDirectory(sourceDirPath string, destinationDirPath string, shouldOverwrite bool) error {
	sourceDirInfo, err := os.Stat(sourceDirPath)
	if err != nil {
		return err
	}

	if !sourceDirInfo.IsDir() {
		return errors.New("Failed to copy directory, '" + sourceDirPath +
			"' is not a directory")
	}

	directoryContents, err := ioutil.ReadDir(sourceDirPath)
	if err != nil {
		return err
	}

	for _, content := range directoryContents {
		contentPath := sourceDirPath + "/" + content.Name()
		if content.IsDir() {
			destinationPath := destinationDirPath + "/" + content.Name()
			err = CopyDirectory(contentPath, destinationPath, shouldOverwrite)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(contentPath, destinationDirPath, shouldOverwrite)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies the source file into the destination path. For example,
// given a source of '/home/user/junk.txt' and a destination of '/tmp/test',
// the resulting file would be '/tmp/test/junk.txt'.
//
// Based on work by "Salvador Dali": https://stackoverflow.com/a/33865286
func CopyFile(sourceFilePath string, destinationDirPath string, shouldOverwrite bool) error {
	sourceFileInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return err
	}

	fullDestinationPath := destinationDirPath + "/" + sourceFileInfo.Name()
	if !shouldOverwrite && Exists(fullDestinationPath) {
		return errors.New("File '" + sourceFileInfo.Name() +
			"' already exists in destination directory '" +
			destinationDirPath + "'")
	}

	parentDirInfo, err := os.Stat(path.Dir(sourceFilePath))
	if err != nil {
		return err
	}

	err = os.MkdirAll(destinationDirPath, parentDirInfo.Mode())
	if err != nil {
		return err
	}

	destFile, err := os.Create(fullDestinationPath)
	if err == nil {
		os.Chmod(fullDestinationPath, sourceFileInfo.Mode())
	} else {
		return err
	}
	defer destFile.Close()

	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

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
// the replacement string already exists, then the function does nothing.
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
