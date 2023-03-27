package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	Model
	TagID      int    `json:"tag_id" gorm:"index"`
	Tag        Tag    `json:"tag"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Content    string `json:"content"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (article *Article) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (article *Article) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}

func ExistArticleByID(id int) bool {
	var article Article
	db.Select("id").Where("id = ?", id).First(&article)
	if article.ID > 0 {
		return true
	}
	return false
}

func GetArticles(pageNum, pageSize int, maps interface{}) (articles []Article) {
	//与其说是pageNum，不如说是articleNum。
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

func GetArticlesTotal(maps interface{}) int {
	var count int64
	db.Model(&Article{}).Where(maps).Count(&count)
	return int(count)
}

func GetArticle(id int) (article Article) {
	db.Where("id = ?", id).Preload("Tag").First(&article)
	return
}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State:     data["state"].(int),
	})
	return true
}

func EditArticle(id int, data map[string]interface{}) bool {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)
	return true
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})
	return true
}

func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Article{})

	return true
}
