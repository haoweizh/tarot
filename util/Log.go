package util

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var socket, info, notice *log.Logger
var socketInfoFile, infoFile, noticeFile *os.File
var socketCount, infoCount, noticeCount int

const logRoot = "./log/"

func initLog(path string) (*log.Logger, *os.File, error) {
	//removeOldFiles()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return nil, nil, err
	}
	return log.New(file, "", log.Ldate|log.Ltime), file, nil
}

//func removeOldFiles() {
//	year, month, date := GetNow().Date()
//	strDate := strconv.Itoa(year) + month.String() + strconv.Itoa(date)
//	err := filepath.Walk(logRoot, func(path string, f os.FileInfo, err error) error {
//		if f == nil {
//			return err
//		}
//		if f.IsDir() {
//			return nil
//		}
//		fmt.Printf(path)
//		if !strings.Contains(f.Name(), strDate) {
//			rmErr := os.Remove(logRoot + f.Name())
//			if rmErr != nil {
//				fmt.Println(logRoot + f.Name() + "can not remove " + rmErr.Error())
//			}
//		}
//		return nil
//	})
//	if err != nil {
//		fmt.Println("can not walk folder " + err.Error())
//	}
//}

func getPath(name string) string {
	year, month, date := GetNow().Date()
	strDate := strconv.Itoa(year) + month.String() + strconv.Itoa(date)
	strTime := strconv.Itoa(GetNow().Hour()) + "_" + strconv.Itoa(GetNow().Minute())
	return logRoot + name + strDate + "_" + strTime + ".log"
}

func SocketInfo(message string) {
	fmt.Println(message)
	if socketCount%10000 == 0 {
		if socketInfoFile != nil {
			socketInfoFile.Close()
		}
		socket, socketInfoFile, _ = initLog(getPath("socketInfo"))
	}
	socket.Println(message)
	socketCount++
}

func Info(message string) {
	//fmt.Println(message)
	if infoCount%10000 == 0 {
		if infoFile != nil {
			infoFile.Close()
		}
		info, infoFile, _ = initLog(getPath("info"))
	}
	info.Println(message)
	infoCount++
}

func Notice(message string) {
	if noticeCount%10000 == 0 {
		if noticeFile != nil {
			noticeFile.Close()
		}
		notice, noticeFile, _ = initLog(getPath("notice"))
	}
	notice.Println(message)
	noticeCount++
}
