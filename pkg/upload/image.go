package upload

import (
	"fmt"
	"gin/pkg/file"
	"gin/pkg/logging"
	"gin/pkg/setting"
	"gin/pkg/util"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// http://127.0.0.1:8000/upload/images/2019/04/04/20190404100000.jpg
func GetImageFullUrl(name string) string {

	return setting.AppSetting.PrefixUrl + "/" + GetImagePath() + name
}

// upload/images/
func GetImagePath() string {

	return setting.AppSetting.ImageSavePath
}

// runtime/upload/images/
func GetImageFullPath() string {

	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	return fileName + ext
}

func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.ImageMaxSize
}

// CheckImage checks if the image exists and has the right permissions
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
