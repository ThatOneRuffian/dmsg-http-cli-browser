package main

import (
	"dmsggui"
	"flag"
)

func main() {

	dmsggui.ClearScreen()
	dmsggui.InitProgramWorkingDir()
	dmsggui.RetryAttemptsUserInput = "3"
	dmsggui.DownloadLocationUserInput = dmsggui.InitDownloadsFolder("")

	//parse user arguments
	flag.StringVar(&dmsggui.DownloadLocationUserInput, "d", dmsggui.DownloadLocationUserInput, "Specify directory to download files to.")
	flag.StringVar(&dmsggui.RetryAttemptsUserInput, "r", dmsggui.RetryAttemptsUserInput, "Specify number of download attempts.")

	flag.Parse()

	//format user provided dir path
	pathByteArray := []byte(dmsggui.DownloadLocationUserInput)
	lastByteChar := pathByteArray[len(pathByteArray)-1]
	const forwardSlash byte = 92
	const backSlash byte = 47

	if lastByteChar == forwardSlash || lastByteChar == backSlash {
		pathByteArray := []byte(dmsggui.DownloadLocationUserInput)
		lastCharDropped := pathByteArray[:len(pathByteArray)-1]
		dmsggui.MainDownloadsLoc = string(lastCharDropped)
	}

	dmsggui.InitDownloadsFolder(dmsggui.DownloadLocationUserInput)

	//attempt to load server cache
	//if config not found then run the first launch wizard
	if !dmsggui.LoadCache() {
		dmsggui.FirstRunWizard()
	}

	for true {
		userChoice := dmsggui.ServerListMainMenu()
		dmsggui.ServerIndexMenuHandler(userChoice)
	}
}
