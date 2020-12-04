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

var DiscoveryServer string = ""

const defaultTerminalWidth = 100

const defaultTerminalHeight = 20

func dmsggetWrapper(publicKey string, downloadLoc string, file string, alternateFileName string, stdOutput bool) {
	var programArgs []string
	downloadInfo := ""
	retryAttempts := RetryAttemptsUserInput
	dmsgDiscFlag := "-dmsg-disc" //work around - dmsgget doesn't default when "" is passed
	DiscoveryServer = stripIllegalChars(DiscoveryServer)

	// change dir back to program working dir for relative paths
	if err := os.Chdir(programCurrentWorkingDir); err != nil {
		fmt.Println("Error navigating program root directory.")
		os.Exit(1)
	}
	//if file contains dirs
	if strings.Contains(file, "/") {
		//extract file name
		fileName := strings.Split(file, "/")
		fileName[0] = fileName[len(fileName)-1]
		alternateFileName = fileName[0]

		downloadInfo = fmt.Sprintf("Downloading %s to %s", string(alternateFileName), downloadLoc)

	} else {
		//if index - better check?
		downloadInfo = fmt.Sprintf("Downloading %s.%s to %s", file, publicKey, downloadLoc)

	}

	//if file exist then overwrite
	if downloadLoc == MainDownloadsLoc {
		os.Remove(downloadLoc + "/" + alternateFileName)
		os.Remove(downloadLoc + "/" + file)
	}

	fetchString := fmt.Sprintf("dmsg://%s:80/%s", publicKey, file)
	stdOutLoc := os.Stdout
	fmt.Println()
	fmt.Println(downloadInfo)
	fmt.Println()
	//open up null to send dmsgget stdout into null for cleaner view
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
		// if not found in path then search current dir for dmsgget

		localDmsgget := fmt.Sprintf(programCurrentWorkingDir + "/dmsgget")
		if runtime.GOOS == "windows" {
			localDmsgget += ".exe"
		}
		dmsggetPath, err = exec.LookPath(localDmsgget)
		if err != nil {
			fmt.Println(fmt.Sprintf("Unable to find dmsgget in PATH or current dir: %s", localDmsgget))
			os.Exit(1)
		}
	}

	programArgs = []string{dmsggetPath, "-t", fmt.Sprint(retryAttempts), "-O", downloadLoc + "/" + alternateFileName}

	if len(DiscoveryServer) > 0 {
		programArgs = append(programArgs, dmsgDiscFlag)
		programArgs = append(programArgs, DiscoveryServer)
	}

	programArgs = append(programArgs, fetchString)

	dmsggetCmd := &exec.Cmd{
		Path:   dmsggetPath,
		Args:   programArgs,
		Stdout: stdOutLoc,
		Stderr: os.Stderr,
	}
	if err := dmsggetCmd.Run(); err != nil {

	}
}

func sttyWrapperGetTerminalHeight() (int, error) {
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
		decodeError := errors.New("Error decoding tput output")
		return returnValue, decodeError
	}

	return returnValue, nil
}

func sttyWrapperGetTerminalWidth() (int, error) {
	returnValue := 0

	cmd := exec.Command("tput", "cols")
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
