package main

import "dmsg_gui"

func main() {
	dmsg_gui.CreateTest()
	/*dmsg_gui.ClearScreen()
	dmsg_gui.InitDownloadsFolder()
	// if config not found then run the first launch wizard
	if !dmsg_gui.LoadCache() {
		dmsg_gui.FirstRunWizard()
	}

	for true {
		userChoice := dmsg_gui.ServerListMainMenu()
		dmsg_gui.ServerIndexMenuHandler(userChoice)
	}*/
}
