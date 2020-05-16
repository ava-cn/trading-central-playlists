package consoles

import (
	"encoding/xml"
	"fmt"
	"github.com/ava-cn/trading-central-playlists/app/models"
	"github.com/ava-cn/trading-central-playlists/databases"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Videos 视频数组
type Videos struct {
	XMLName xml.Name `xml:"videos"`
	Video   []Video  `xml:"video"`
}

// Video 视频
type Video struct {
	Version      string    `xml:"version"`
	ID           uint64    `xml:"id"`
	Title        string    `xml:"title"`
	Description  string    `xml:"description"`
	CreatedAt    time.Time `xml:"created_at"`
	URL          string    `xml:"url"`
	ThumbnailURL string    `xml:"thumbnail_url"`
	ImageURL     string    `xml:"image_url"`
	PlayCount    int       `xml:"play_count"`
	Duration     float64   `xml:"duration"`
	Height       int       `xml:"height"`
	Width        int       `xml:"width"`
}

// 定义一个通道存储视频数据
var CurrentVideoListFromXMLChan = make(chan Video, 10)

func Run() {
	// 获取数据并存入Chan
	go FetchFormURL()

	// 存储数据到数据库
	go StoreToDatabase()

}

// 发送请求获取资源存储到videoListChan中
func FetchFormURL() {
	var (
		data     []byte
		response *http.Response
		err      error
		videos   Videos
		video    Video
	)

	if response, err = http.Get(viper.GetString("app.xml_url")); err != nil {
		log.Fatalf("failed to fetch remote url, err: %s", err.Error())
		return
	}
	defer response.Body.Close()

	if data, err = ioutil.ReadAll(response.Body); err != nil {
		fmt.Printf("open file failed, err: %s", err.Error())
		return
	}

	if err = xml.Unmarshal(data, &videos); err != nil {
		fmt.Printf("XML file unmasrshaler fialed, err: %s", err.Error())
		return
	}

	for _, video = range videos.Video {
		// 查询对应的视频是存在于数据库记录中，如果存在则记录，如果不存在则记录
		if !models.IsVideoExists(databases.GetDB(), video.ID) {
			CurrentVideoListFromXMLChan <- video
		}

	}
}

func StoreToDatabase() {
	var (
		video         models.Videos
		videoFromChan Video
	)

	for {
		select {
		case videoFromChan = <-CurrentVideoListFromXMLChan:
			video = models.Videos{
				VideoID:                 videoFromChan.ID,
				VideoTitle:              videoFromChan.Title,
				VideoCreatedAt:          models.Time(videoFromChan.CreatedAt),
				VideoDuration:           videoFromChan.Duration,
				VideoWidth:              videoFromChan.Width,
				VideoHeight:             videoFromChan.Height,
				OriginVideoUrl:          videoFromChan.URL,
				OriginVideoThumbnailUrl: videoFromChan.ThumbnailURL,
				OriginVideoImageUrl:     videoFromChan.ImageURL,
				VideoExtras:             nil,
			}

			databases.GetDB().Create(&video)

			// 将视频和图片资源上传到七牛云
			go StoreToStorage(&video)
		}
	}
}

func StoreToStorage(video *models.Videos) {
	// 存储视频

	// 存储图片

}
