package models

import (
	"github.com/jinzhu/gorm"
)

type Videos struct {
	ID uint `json:"id" gorm:"primary_key;comment:'主键id'"`

	VideoID        uint64  `json:"video_id" gorm:"type:bigint unsigned;not null;unique;comment:'视频ID'"`
	VideoTitle     string  `json:"video_title" gorm:"size:255;not null;comment:'视频名称'"`
	VideoCreatedAt Time    `json:"video_created_at" gorm:"type:datetime;comment:'视频创建时间'"`
	VideoDuration  float64 `json:"video_duration" gorm:"type:decimal(5,2);comment:'视频时间'" sql:"not null;default 0.0;"`
	VideoWidth     int     `json:"video_width" gorm:"column:video_width;comment:'视频宽度'" sql:"not null; default:0; type:int unsigned"`
	VideoHeight    int     `json:"video_height" gorm:"column:video_height;comment:'视频高度'" sql:"not null; default:0; type:int unsigned"`

	OriginVideoUrl          string `json:"origin_video_url" gorm:"size:255;not null;comment:'原始视频地址'"`
	OriginVideoThumbnailUrl string `json:"origin_video_thumbnail_url" gorm:"size:255;not null;comment:'原始视频缩略图地址'"`
	OriginVideoImageUrl     string `json:"origin_video_image_url" gorm:"size:255;not null;comment:'原始视频图片地址'"`

	VideoUrl          string `json:"video_url" gorm:"size:255;comment:'视频地址'"`
	VideoThumbnailUrl string `json:"video_thumbnail_url" gorm:"size:255;comment:'视频缩略图地址'"`
	VideoImageUrl     string `json:"video_image_url" gorm:"size:255;comment:'视频图片地址'"`

	VideoExtras JSON `json:"video_extras,omitempty" gorm:"column:video_extras;comment:'视频额外数据'" sql:"type:json"`

	CreatedAt Time  `json:"created_at" gorm:"type:datetime;comment:'创建时间'"`
	UpdatedAt Time  `json:"updated_at" gorm:"type:datetime;comment:'更新时间'"`
	DeletedAt *Time `json:"deleted_at" gorm:"type:datetime;comment:'删除时间'" sql:"index"`
}

// 通过视频ID查找视频是否存在
func IsVideoExists(db *gorm.DB, VideoID uint64) bool {
	var video Videos

	db.Where("video_id = ?", VideoID).First(&video)

	if video.VideoID != 0 {
		return true
	}

	return false
}
