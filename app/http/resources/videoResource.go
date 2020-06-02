package resources

import (
	"github.com/ava-cn/trading-central-playlists/app/consoles"
	"github.com/ava-cn/trading-central-playlists/app/models"
)

type VideosResource struct {
	VideoID uint64 `json:"video_id"`

	VideoTitle        string      `json:"video_title"`
	VideoCreatedAt    models.Time `json:"video_created_at"`
	VideoDuration     float64     `json:"video_duration"`
	VideoWidth        int         `json:"video_width"`
	VideoHeight       int         `json:"video_height"`
	VideoUrl          string      `json:"video_url"`
	VideoThumbnailUrl string      `json:"video_thumbnail_url"`
	VideoImageUrl     string      `json:"video_image_url"`

	CreatedAt models.Time `json:"created_at"`
	UpdatedAt models.Time `json:"updated_at"`
}

func VideoCollection(videos []*models.Videos) (videoResources []VideosResource) {
	var video *models.Videos

	for _, video = range videos {
		videoResources = append(videoResources, *VideoShow(video))
	}

	return

}

// 视频详情
func VideoShow(video *models.Videos) *VideosResource {
	return &VideosResource{
		VideoID:           video.VideoID,
		VideoTitle:        video.VideoTitle,
		VideoCreatedAt:    video.VideoCreatedAt,
		VideoDuration:     video.VideoDuration,
		VideoWidth:        video.VideoWidth,
		VideoHeight:       video.VideoHeight,
		VideoUrl:          consoles.GetFile(video.VideoUrl),
		VideoThumbnailUrl: consoles.GetFile(video.VideoThumbnailUrl),
		VideoImageUrl:     consoles.GetFile(video.VideoImageUrl),
		CreatedAt:         video.CreatedAt,
		UpdatedAt:         video.UpdatedAt,
	}
}
