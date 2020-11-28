package main

import (
	"dmsggui"
	"flag"
)

func main() {

	dmsggui.ClearScreen()
	dmsggui.DownloadLocationUserInput = dmsggui.InitDownloadsFolder()
	dmsggui.RetryAttemptsUserInput = "3"
	// if config not found then run the first launch wizard

	// parse user arguments
	flag.StringVar(&dmsggui.DownloadLocationUserInput, "d", dmsggui.DownloadLocationUserInput, "Specify directory to download files to.")
	flag.StringVar(&dmsggui.RetryAttemptsUserInput, "r", dmsggui.RetryAttemptsUserInput, "Specify number of download attempts.")

	flag.Parse()

	if !dmsggui.LoadCache() {
		dmsggui.FirstRunWizard()
	}

	for true {
		userChoice := dmsggui.ServerListMainMenu()
		dmsggui.ServerIndexMenuHandler(userChoice)
	}
}
