package dmsggui

import "fmt"

var downloadQueue = make(map[int]downloadItem)

type downloadItem struct {
	fileName           string
	fileSize           float64
	downloadStatus     *bool
	serverFriendlyName string
}

//InitMuxDownload adds items to the multiplexed download queue
func initMuxDownload(serverPublicKey string, MainDownloadsLoc string, _fileName string) {
	//create object to track file stats and add object to queue
	newIndex := len(downloadQueue)
	boolVar := true
	newDownloadItem := downloadItem{
		fileName:           _fileName,
		fileSize:           navPtr.files[_fileName],
		downloadStatus:     &boolVar,
		serverFriendlyName: serverPublicKey,
	}

	downloadQueue[newIndex] = newDownloadItem

	//start download
	go dmsggetWrapper(serverPublicKey, MainDownloadsLoc, getPresentWorkingDirectory()+_fileName, "", false)

}

func clearFinishedDownloadsFromQueue() {
	tmpList := make(map[int]downloadItem)

	for index := range downloadQueue {
		fmt.Println("index", index)
		if *downloadQueue[index].downloadStatus {
			tmpList[index] = downloadQueue[index]
		}
	}
	downloadQueue = tmpList
}
