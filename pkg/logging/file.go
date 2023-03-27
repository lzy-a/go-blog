package logging

import (
	"fmt"
	"gin/pkg/file"
	"gin/pkg/setting"
	"os"
	"time"
)

var (
	LogSavePath     = setting.AppSetting.LogSavePath
	LogSaveName     = setting.AppSetting.LogSaveName
	LogFileExt      = setting.AppSetting.LogFileExt
	TimeFormat      = setting.AppSetting.TimeFormat
	RuntimeRootPath = setting.AppSetting.RuntimeRootPath
)

func getLogFilePath() string {
	return RuntimeRootPath + LogSavePath
}

func getLogFileName() string {
	return LogSaveName + time.Now().Format(TimeFormat) + "." + LogFileExt
}

func openLogFile(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}
	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}
	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("Fail to OpenFile :%v", err)
	}
	return f, nil
}
