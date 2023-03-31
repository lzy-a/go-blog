package export

import "gin/pkg/setting"

func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + getExcelPath() + name
}

func getExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + getExcelPath()
}
