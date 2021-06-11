package consoles

import (
    "encoding/json"
    "encoding/xml"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/ava-cn/trading-central-playlists/app/models"
    "github.com/ava-cn/trading-central-playlists/app/supports"
    "github.com/ava-cn/trading-central-playlists/databases"
    "github.com/spf13/viper"
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

// RunSync 文件同步
func RunSync() {
    // 获取数据并存入Chan
    log.Println("FetchFormURL start running...")
    go FetchFormURL()

    // 检查数据库未同步的数据
    log.Println("CheckSyncedStatus start running...")
    go CheckSyncedStatus()
}

// FetchFormURL 发送请求获取资源存储到videoListChan中
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
            // 存储数据到数据库
            log.Println("StoreToDatabase start running...")
            StoreToDatabase(video)
        }
    }
}

// StoreToDatabase 存储到数据库
func StoreToDatabase(video Video) {
    var (
        videoModel *models.Videos

        VideoExtras     models.VideoExtras
        videoExtrasJson []byte
    )

    // 获取最终的URL地址
    VideoExtras.RedirectVideoURL, _ = supports.GetRedirectURL(video.URL)
    VideoExtras.RedirectVideoImageURL, _ = supports.GetRedirectURL(video.ImageURL)
    VideoExtras.RedirectVideoThumbnailURL, _ = supports.GetRedirectURL(video.ThumbnailURL)

    // 获取文件名
    VideoExtras.RealVideoName, _ = supports.GetFileNameFromURL(VideoExtras.RedirectVideoURL)
    VideoExtras.RealVideoImageName, _ = supports.GetFileNameFromURL(VideoExtras.RedirectVideoImageURL)
    VideoExtras.RealVideoThumbnailName, _ = supports.GetFileNameFromURL(VideoExtras.RedirectVideoThumbnailURL)

    videoExtrasJson, _ = json.Marshal(VideoExtras)

    videoModel = &models.Videos{
        VideoID:                 video.ID,
        VideoTitle:              video.Title,
        VideoCreatedAt:          models.Time(video.CreatedAt),
        VideoDuration:           video.Duration,
        VideoWidth:              video.Width,
        VideoHeight:             video.Height,
        OriginVideoUrl:          video.URL,
        OriginVideoThumbnailUrl: video.ThumbnailURL,
        OriginVideoImageUrl:     video.ImageURL,
        VideoExtras:             videoExtrasJson,
        Synced:                  false,
    }

    if !models.IsVideoExists(databases.GetDB(), videoModel.VideoID) {
        databases.GetDB().Create(videoModel)

        // 将视频和图片资源上传到七牛云
        go StoreToStorage(*videoModel)
    } else {
        log.Printf("视频%d已经存在\n", videoModel.VideoID)
    }
}

// CheckSyncedStatus 检查未同步的数据 synced = 0
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
    if viper.GetBool("qiniu.enabled") {
        log.Printf("we have %d tasks should store to qiniu storage", unSyncedCount)
    }

    if unSyncedCount >= 0 {
        for _, video = range videos {
            // 将视频和图片资源上传到七牛云
            go StoreToStorage(video)
        }
    }
}

// StoreToStorage 保存到七牛云存储
func StoreToStorage(video models.Videos) {

    if !viper.GetBool("qiniu.enabled") {
        log.Println("ignore storage to qiniu....")
        return
    } else {
        log.Printf("start storage to qiniu: %#v", video.OriginVideoUrl)
    }

    var (
        videoPathPrefix   = "videos/" + strconv.Itoa(int(video.VideoID)) + "/"
        imagePathPrefix   = "images/" + strconv.Itoa(int(video.VideoID)) + "/"
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
    databases.GetDB().Model(video).Update(models.Videos{
        VideoUrl:          videoKey,
        VideoThumbnailUrl: imageKey,
        VideoImageUrl:     thumbnailImageKey,
        Synced:            true,
    })

ERR:
    log.Println("uploadToQiniu: ", err)
}
