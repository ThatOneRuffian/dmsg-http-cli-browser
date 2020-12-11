package dmsggui

import (
	"fmt"
	"sort"
)

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
	sortedIndexes := []int{}
	newIndex := 0
	for index := range downloadQueue {
		sortedIndexes = append(sortedIndexes, index)
	}
	sort.Ints(sortedIndexes)
	for index := range sortedIndexes {
		fmt.Println("index", index)
		if *downloadQueue[index].downloadStatus {
			tmpList[newIndex] = downloadQueue[index]
			newIndex++
		}
	}
	downloadQueue = tmpList
}
