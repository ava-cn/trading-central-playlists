## Trading Central Playlists

## 说明

从 [Trading Central PlayLists](https://video.tradingcentral.com/playlists/23125.xml) 获取对应的视频和图片资源存储到七牛云存储，并提供统一的数据接口。

## 下载

```
git clone git@github.com:ava-cn/trading-central-playlists.git
```

## 配置

配置文件在`configs/app.yml`

```
server:
 port: 8080

database:
 driverName: mysql
 host: 127.0.0.1
 port: 3306
 database: trading_central_playlists
 user: root
 password:
 charset: utf8
 local: Asia/Shanghai

app:
 xml_url: https://video.tradingcentral.com/playlists/23125.xml

qiniu:
 bucket:
 accessKey:
 secretKey:
 # 空间对应的机房
 useHTTPS: true
 useCdnDomains: true
```


## 接口

暂未开发
