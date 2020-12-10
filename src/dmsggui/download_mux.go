package dmsggui

var downloadQueue map[int]downloadItem

type downloadItem struct {
	fileName           string
	fileSize           float64
	downloadStatus     float64
	serverFriendlyName string
}

//InitMuxDownload adds items to the multiplexed download queue
func initMuxDownload(serverPublicKey string, MainDownloadsLoc string, fileName string) {
	//create object to track file stats and add object to queue

	//start download
	go dmsggetWrapper(serverPublicKey, MainDownloadsLoc, getPresentWorkingDirectory()+fileName, "", false)

}
