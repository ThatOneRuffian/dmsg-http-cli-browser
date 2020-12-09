package main

import (
	"dmsggui"
	"flag"
)

func main() {

	dmsggui.ClearScreen()
	dmsggui.InitProgramWorkingDir()
	dmsggui.RetryAttemptsUserInput = "3"

	//parse user arguments
	flag.StringVar(&dmsggui.DiscoveryServer, "dmsg-disc", "", "Specify the dmsg discovery URL. Default is dmsgget's default discovery URL.")
	flag.StringVar(&dmsggui.DownloadLocationUserInput, "d", dmsggui.DownloadLocationUserInput, "Specify directory to download files to.")
	flag.StringVar(&dmsggui.RetryAttemptsUserInput, "r", dmsggui.RetryAttemptsUserInput, "Specify number of download attempts.")
	flag.StringVar(&dmsggui.DiscoveryServerPort, "p", "80", "Specify the port number of the discovery server.")

	flag.Parse()

	//format user provided dir path
	pathByteArray := []byte(dmsggui.DownloadLocationUserInput)
	const forwardSlash byte = 92
	const backSlash byte = 47
	var lastCharByte byte

	if len(dmsggui.DownloadLocationUserInput) > 0 {
		lastCharByte = pathByteArray[len(pathByteArray)-1]
		if lastCharByte == forwardSlash || lastCharByte == backSlash {
			downloadPathRemovedDirSymbol := pathByteArray[:len(pathByteArray)-1]
			dmsggui.MainDownloadsLoc = string(downloadPathRemovedDirSymbol)
		}
	}

	dmsggui.InitDownloadsFolder(dmsggui.DownloadLocationUserInput)

	//attempt to load server cache
	//if cache not found then run the first launch wizard
	if !dmsggui.LoadCache() {
		dmsggui.FirstRunWizard()
	}

	for true {
		userChoice := dmsggui.ServerListMainMenu()
		dmsggui.ServerIndexMenuHandler(userChoice)
	}
}
