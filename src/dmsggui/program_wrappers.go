package dmsggui

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const defaultTerminalWidth = 100

const defaultTerminalHeight = 35

func dmsggetWrapper(publicKey string, downloadLoc string, file string, alternateFileName string, stdOutput bool) {
	downloadInfo := ""
	retryAttempts := RetryAttemptsUserInput
	if strings.Contains(file, "/") {
		fileName := strings.Split(file, "/")

		fileName[0] = fileName[len(fileName)-1]
		alternateFileName = fileName[0]
		downloadInfo = fmt.Sprintf("Downloading %s to %s", string(alternateFileName), downloadLoc)

	} else {
		downloadInfo = fmt.Sprintf("Downloading %s.%s to %s", file, publicKey, downloadLoc)

	}
	if downloadLoc == mainDownloadsLoc {
		clearFile(downloadLoc + "/" + alternateFileName)
		clearFile(downloadLoc + "/" + file)
	}

	fetchString := fmt.Sprintf("dmsg://%s:80/%s", publicKey, file)
	stdOutLoc := os.Stdout
	fmt.Println()
	fmt.Println(downloadInfo)
	fmt.Println()

	if !stdOutput {
		nullWriteLocation := ""
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			nullWriteLocation = "/dev/null"
		} else if runtime.GOOS == "windows" {
			nullWriteLocation = "./nul"
		}

		nullFile, err := os.OpenFile(nullWriteLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer nullFile.Close()
		if err != nil {
			fmt.Println("Error opening null for writing")
		}
		stdOutLoc = nullFile
	}
	dmsggetPath, err := exec.LookPath("dmsgget")
	if err != nil {
		ClearScreen()
		log.Fatal("Unable to find dmsgget. Make sure dmsgget's location is included in your PATH variable.")
	}

	dmsggetCmd := &exec.Cmd{
		Path:   dmsggetPath,
		Args:   []string{dmsggetPath, "-t", fmt.Sprint(retryAttempts), "-O", downloadLoc + "/" + alternateFileName, fetchString},
		Stdout: stdOutLoc,
		Stderr: os.Stderr,
	}
	if err := dmsggetCmd.Run(); err != nil {

	}
}

func sttyWrapperGetTerminalHeight() (int, error) {
	returnValue := 0

	cmd := exec.Command("tput", "lines")
	fmt.Println(cmd.Args)
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
		decodeError := errors.New("Error decoding tput output")
		return returnValue, decodeError
	}

	return returnValue, nil
}

func sttyWrapperGetTerminalWidth() (int, error) {
	returnValue := 0

	cmd := exec.Command("tput", "cols")
	fmt.Println(cmd.Args)
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
		decodeError := errors.New("Error decoding tput output")
		return returnValue, decodeError
	}

	return returnValue, nil
}

func getTerminalDims(bufferHeight int) (int, int) {
	terminalHeightAvailable := 1
	terminalWidthAvailable := 1
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		terminalHeight, heightError := sttyWrapperGetTerminalHeight()
		terminalWidth, widthError := sttyWrapperGetTerminalWidth()
		if heightError != nil || widthError != nil {
			fmt.Println("Error fetching terminal dimensions")
			fmt.Println(heightError)
			fmt.Println(widthError)
			terminalHeightAvailable = defaultTerminalHeight //default on error
			terminalWidthAvailable = defaultTerminalWidth

		} else {
			terminalHeight -= bufferHeight
			terminalHeightAvailable = terminalHeight
			terminalWidthAvailable = terminalWidth
		}
	} else if runtime.GOOS == "windows" {
		//todo auto sizing of output
		terminalHeightAvailable = defaultTerminalHeight //default on error
		terminalWidthAvailable = defaultTerminalWidth
	}
	return terminalHeightAvailable, terminalWidthAvailable
}
