package tagservice

import (
	"encoding/json"
	"gin/models"
	"gin/pkg/export"
	"gin/pkg/gredis"
	"gin/pkg/logging"
	cacheservice "gin/service/cache_service"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int

	PageNum  int
	PageSize int
}

func (t *Tag) ExistByID() (bool, error) {
	return models.ExistTagByID(t.ID)
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var cacheTags, tags []models.Tag
	cacheservice := cacheservice.Tag{
		State:    t.State,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cacheservice.GetTagsKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}
	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, tags, 3600)
	return tags, nil
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Delete() error {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}
	err := models.EditTag(t.ID, data)
	if err != nil {
		return err
	}
	cacheservice := cacheservice.Tag{ID: t.ID}
	key := cacheservice.GetTagsKey()
	tag, err := models.GetTag(t.ID)
	if err != nil {
		return err
	}
	gredis.Set(key, tag, 3600)
	return nil
}

func (t *Tag) Export() (string, error) {
	tags, err := t.GetAll()
	// fmt.Println(tags)
	if err != nil {
		return "", err
	}
	file := excelize.NewFile()
	index := file.NewSheet("标签信息")

	titles := []string{"ID", "名称", "创建人", "创建时间", "修改人", "修改时间"}
	for col, v := range titles {
		file.SetCellValue("标签信息", excelize.ToAlphaString(col)+"1", v)
	}
	for row, v := range tags {
		values := []interface{}{
			v.ID, v.Name, v.CreatedBy, v.CreatedOn, v.ModifiedBy, v.ModifiedOn,
		}
		for col, value := range values {
			file.SetCellValue("Sheet1", excelize.ToAlphaString(col)+strconv.Itoa(row+2), value)
		}
	}

	time := strconv.Itoa(int(time.Now().Unix()))
	filename := "tags-" + time + ".xlsx"
	fullPath := export.GetExcelFullPath() + filename
	if err := file.SaveAs(fullPath); err != nil {
		return "", err
	}
	file.SetActiveSheet(index)
	return filename, nil
}

func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}
