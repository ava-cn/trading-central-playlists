## Trading Central Playlists

## 说明

从 [Trading Central PlayLists](https://video.tradingcentral.com/playlists/23125.xml) 获取对应的视频和图片资源存储到七牛云存储，并提供统一的数据接口。

## 下载

```
git clone git@github.com:ava-cn/trading-central-playlists.git
```

## 编译和执行

- 编译
```
make build_linux_amd64
```

- 运行
```
./trading-central-playlists -c configs/app.yml # 通过自定义配置文件运行服务
```

> 运行之前需要保证MySQL、Qiniu等依赖服务配置正常连接。

## Docker

### 构建docker镜像

```
docker build -t curder/trading-central-playlists .
```

### 启动容器

- 依赖MySQL服务

```
docker pull mysql:5.7.29 # 拉取镜像

docker run --name trading-central-playlists-mysql \
    -p 33068:3306 \
    -v ~/.docker/trading-central-playlists-mysql/data:/var/lib/mysql \
    -e MYSQL_ROOT_PASSWORD=root \
    -e MYSQL_DATABASE=trading_central_playlists \
    -e MYSQL_USER=trading_central_playlists \
    -e MYSQL_PASSWORD=trading_central_playlists \
    -d mysql:5.7.29 \
    --character-set-server=utf8mb4 \
    --collation-server=utf8mb4_unicode_ci
```

- 启动项目容器
```
cp app.yml app.production.yml # 拷贝并修改项目配置
docker run --name trading-central-playlists \
    --link trading-central-playlists-mysql:mysql \
    -v app.yml:$HOME/.trading-central-playlists/app.yml \
    -p 8088:80 \
    curder/trading-central-playlists \
    -d
```

### 数据库备份和还原

- 备份
```
docker exec -it trading-central-playlists-mysql mysqldump -utrading_central_playlists -ptrading_central_playlists trading_central_playlists| gzip > trading_central_playlists.tar.gz
```

- 还原
```
gzip -d trading_central_playlists.tar.gz # 解压缩，原.tar.gz文件会被删除
mysql -utrading_central_playlists -p trading_central_playlists < trading_central_playlists.tar
```

### 测试服务

```
curl http://127.0.0.1:8088/ping
```

## docker-compose

### 拷贝并修改配置文件

```
cp app.yml app.production.yml
```

### 启动

```
docker-compose up -d
```

### CentOS services

```
sudo cp docker-compose-trading-central-playlists.service /usr/lib/systemd/system/.
```

- 查看状态
`systemctl status docker-compose-trading-central-playlists`

- 启动服务
```
systemctl start docker-compose-trading-central-playlists
```

- 暂停服务
```
systemctl stop docker-compose-trading-central-playlists
``` 

## 接口

- **GET** `/ping`
    返回值 `{"message": "PONG"}`

- **GET** `/playlists`
    返回值 `{"code": 200, "data": [], "message": "数据获取成功"}`
