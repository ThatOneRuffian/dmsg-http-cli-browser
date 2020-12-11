package dmsggui

var downloadQueue = make(map[int]downloadItem)

type downloadItem struct {
	fileName           string
	fileSize           float64
	downloadStatus     float64
	serverFriendlyName string
}

//InitMuxDownload adds items to the multiplexed download queue
func initMuxDownload(serverPublicKey string, MainDownloadsLoc string, _fileName string) {
	//create object to track file stats and add object to queue
	newIndex := len(downloadQueue)

	newDownloadItem := downloadItem{
		fileName:           _fileName,
		fileSize:           navPtr.files[_fileName],
		serverFriendlyName: serverPublicKey,
	}

	downloadQueue[newIndex] = newDownloadItem

	//start download
	go dmsggetWrapper(serverPublicKey, MainDownloadsLoc, getPresentWorkingDirectory()+_fileName, "", false)

}
