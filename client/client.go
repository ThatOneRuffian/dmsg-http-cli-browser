package main

import "dmsggui"

func main() {

	dmsggui.ClearScreen()
	dmsggui.InitDownloadsFolder()
	// if config not found then run the first launch wizard
	if !dmsggui.LoadCache() {
		dmsggui.FirstRunWizard()
	}

	for true {
		userChoice := dmsggui.ServerListMainMenu()
		dmsggui.ServerIndexMenuHandler(userChoice)
	}
}
