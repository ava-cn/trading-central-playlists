package consoles

import (
	"encoding/json"
	"encoding/xml"
	"github.com/ava-cn/trading-central-playlists/app/models"
	"github.com/ava-cn/trading-central-playlists/app/supports"
	"github.com/ava-cn/trading-central-playlists/databases"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// 文件同步
func RunSync() {
	// 获取数据并存入Chan
	log.Println("FetchFormURL start running...")
	go FetchFormURL()

	// 存储数据到数据库
	log.Println("StoreToDatabase start running...")
	go StoreToDatabase()

	// 检查数据库未同步的数据
	log.Println("CheckSyncedStatus start running...")
	go CheckSyncedStatus()
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
		log.Printf("open file failed, err: %s", err.Error())
		return
	}

	if err = xml.Unmarshal(data, &videos); err != nil {
		log.Printf("XML file unmasrshaler fialed, err: %s", err.Error())
		return
	}

	for _, video = range videos.Video {
		// 查询对应的视频是存在于数据库记录中，如果存在则记录，如果不存在则记录
		if !models.IsVideoExists(databases.GetDB(), video.ID) {

			CurrentVideoListFromXMLChan <- video
		}
	}
}

// 存储到数据库
func StoreToDatabase() {
	var (
		video         *models.Videos
		videoFromChan Video

		VideoExtras     models.VideoExtras
		videoExtrasJson []byte
	)

	for {
		select {
		case videoFromChan = <-CurrentVideoListFromXMLChan:

			time.Sleep(1 * time.Second) // 休眠1S

			// 获取最终的URL地址
			VideoExtras.RedirectVideoURL, _ = supports.GetRedirectURL(videoFromChan.URL)
			VideoExtras.RedirectVideoImageURL, _ = supports.GetRedirectURL(videoFromChan.ImageURL)
			VideoExtras.RedirectVideoThumbnailURL, _ = supports.GetRedirectURL(videoFromChan.ThumbnailURL)

			// 获取文件名
			VideoExtras.RealVideoName, _ = supports.GetFileNameFromURL(VideoExtras.RedirectVideoURL)
			VideoExtras.RealVideoImageName, _ = supports.GetFileNameFromURL(VideoExtras.RedirectVideoImageURL)
			VideoExtras.RealVideoThumbnailName, _ = supports.GetFileNameFromURL(VideoExtras.RedirectVideoThumbnailURL)

			videoExtrasJson, _ = json.Marshal(VideoExtras)

			video = &models.Videos{
				VideoID:                 videoFromChan.ID,
				VideoTitle:              videoFromChan.Title,
				VideoCreatedAt:          models.Time(videoFromChan.CreatedAt),
				VideoDuration:           videoFromChan.Duration,
				VideoWidth:              videoFromChan.Width,
				VideoHeight:             videoFromChan.Height,
				OriginVideoUrl:          videoFromChan.URL,
				OriginVideoThumbnailUrl: videoFromChan.ThumbnailURL,
				OriginVideoImageUrl:     videoFromChan.ImageURL,
				VideoExtras:             videoExtrasJson,
				Synced:                  false,
			}

			if !models.IsVideoExists(databases.GetDB(), video.VideoID) {
				databases.GetDB().Create(video)

				// 将视频和图片资源上传到七牛云
				go StoreToStorage(video)
			} else {
				log.Printf("视频%d已经存在\n", video.VideoID)
			}

		default:
			time.Sleep(1 * time.Second) // 休眠1S
		}
	}
}

// 检查未同步的数据 synced = 0
func CheckSyncedStatus() {
	var (
		video         models.Videos
		videos        []models.Videos
		unSyncedCount int
		err           error
	)

	if err = databases.GetDB().Where("synced = ?", false).Find(&videos).Count(&unSyncedCount).Error; err != nil {
		log.Println(err)
		return
	}

	log.Println("CheckSyncedStatus func running...")
	log.Printf("we have %d tasks should sync", unSyncedCount)

	if unSyncedCount >= 0 {
		for _, video = range videos {
			// 将视频和图片资源上传到七牛云
			go StoreToStorage(&video)
		}
	}
}

// 保存到七牛云存储
func StoreToStorage(video *models.Videos) {
	var (
		videoPathPrefix   = "trading-central/videos/" + strconv.Itoa(int(video.VideoID)) + "/"
		imagePathPrefix   = "trading-central/images/" + strconv.Itoa(int(video.VideoID)) + "/"
		videoKey          string
		imageKey          string
		thumbnailImageKey string
		videoExtras       models.VideoExtras
		err               error
	)

	_ = json.Unmarshal(video.VideoExtras, &videoExtras)

	// 获取文件名
	if videoKey, err = UploadToQiniu(video.OriginVideoUrl, videoPathPrefix+videoExtras.RealVideoName); err != nil {
		goto ERR
	}
	// 存储图片
	if imageKey, err = UploadToQiniu(video.OriginVideoImageUrl, imagePathPrefix+videoExtras.RealVideoImageName); err != nil {
		goto ERR
	}

	if thumbnailImageKey, err = UploadToQiniu(video.OriginVideoThumbnailUrl, imagePathPrefix+videoExtras.RealVideoThumbnailName); err != nil {
		goto ERR
	}

	// 存储到数据库
	databases.DB.Model(video).Update(models.Videos{
		VideoUrl:          videoKey,
		VideoThumbnailUrl: imageKey,
		VideoImageUrl:     thumbnailImageKey,
		Synced:            true,
	})

ERR:
}
