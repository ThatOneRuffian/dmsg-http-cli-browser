package dmsg_gui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func dmsggetWrapper(publicKey string, downloadLoc string, file string, alternateFileName string, stdOutput bool) bool {
	downloadInfo := ""

	if strings.Contains(file, "/") {
		fileName := strings.Split(file, "/")

		fileName[0] = fileName[len(fileName)-1]
		alternateFileName = fileName[0]
		downloadInfo = fmt.Sprintf("Downloading %s to %s", string(alternateFileName), downloadLoc)

	} else {
		downloadInfo = fmt.Sprintf("Downloading %s to %s", file, downloadLoc)

	}
	if downloadLoc == MainDownloadsLoc {
		ClearFile(downloadLoc + "/" + alternateFileName)
		ClearFile(downloadLoc + "/" + file)
	}

	fmt.Println(alternateFileName)
	fetchString := fmt.Sprintf("dmsg://%s:80/%s", publicKey, file)
	returnValue := true
	stdOutLoc := os.Stdout
	fmt.Println(downloadInfo)
	if !stdOutput {
		nullFile, err := os.OpenFile("/dev/null", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening /dev/null for writing")
		}
		defer nullFile.Close()
		stdOutLoc = nullFile
	}
	dmsggetPath, err := exec.LookPath("dmsgget")
	if err != nil {
		fmt.Println(err)
	}

	dmsggetCmd := &exec.Cmd{
		Path:   dmsggetPath,
		Args:   []string{dmsggetPath, "-O", downloadLoc + "/" + alternateFileName, fetchString},
		Stdout: stdOutLoc,
		Stderr: os.Stderr,
	}
	if err := dmsggetCmd.Run(); err != nil {
		fmt.Println("There was an error fetching the file", err)
		// file exists?
		returnValue = false
	}
	return returnValue
}

func SttyWrapperGetTerminalHeight() (int, error) {
	returnValue := 0
	cmd := exec.Command("tput", "lines")
	cmd.Stderr = os.Stderr

	stdOut, err := cmd.StdoutPipe()
	if nil != err {
		log.Fatalf("Error attaching to tput stdout: %s", err.Error())
	}
	stdOutReader := bufio.NewReader(stdOut)
	go func(stdOutReader io.Reader) {
		scanner := bufio.NewScanner(stdOutReader)
		for scanner.Scan() {
			returnValue, err = strconv.Atoi(scanner.Text())
			if err != nil {
				returnValue = 0
			}
		}
	}(stdOutReader)

	if err := cmd.Start(); nil != err {
		fmt.Println(fmt.Sprintf("Error starting program: %s, %s", cmd.Path, err.Error()))
		fmt.Println("Make sure that tput is installed on your system")
	}
	cmd.Wait()
	if returnValue == 0 {
		docdeError := errors.New("Error decoding tput output")
		return returnValue, docdeError
	} else {
		return returnValue, nil
	}

}
