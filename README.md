# AVMeta

![Build](https://github.com/xnervwang/AVMeta/workflows/Build/badge.svg)
![Release](https://github.com/xnervwang/AVMeta/workflows/Release/badge.svg)
[![codecov](https://codecov.io/gh/xnervwang/AVMeta/branch/master/graph/badge.svg)](https://codecov.io/gh/xnervwang/AVMeta)
[![Go Report Card](https://goreportcard.com/badge/github.com/xnervwang/AVMeta)](https://goreportcard.com/report/github.com/xnervwang/AVMeta)
![GitHub](https://img.shields.io/github/license/xnervwang/AVMeta)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xnervwang/AVMeta)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/xnervwang/AVMeta)
![GitHub All Releases](https://img.shields.io/github/downloads/xnervwang/AVMeta/total)

AV 元数据刮削器，使用 Golang 语言编写，具有多线程、全兼容等特点。

通过文件名称自动计算影片番号，并访问各官网或 Jav 类网站获取元数据信息。

获取到元数据后，自动下载并剪切封面图片，并按照指定路径存储电影、元数据、封面。

## 目录

* [FAQ](#FAQ)
* [编译](#编译)
* [配置](#配置)
* [使用](#使用)
    * [头像](#头像)
        * [本地下载](#本地下载)
        * [本地入库](#本地入库)
    * [刮削](#刮削)
        * [NFO刮削](#NFO刮削)
        * [群晖刮削](#群晖刮削)
    * [转换](#转换)
* [鸣谢](#鸣谢)

## FAQ

1. 什么是元数据？
> 元数据就是电影的详细信息，包含：封面、简介、演员、标题等……

2. AVMeta 有什么用？
> 方便整理AV电影而已。

3. 为什么我使用不了？
> 请以下方格式将错误信息填写到 [issue](https://github.com/xnervwang/AVMeta/issues/new) 中。

```bash
操作系统： Windows 7 x64
Go版本： 1.13
AVMeta版本： v1.0.0
配置信息：
将敏感信息替换为*号
错误信息：
文件/番号: [xxx.mp4/xxx] 刮削失败, 错误原因: xxx
```

## 编译

不想编译，可直接在 [发布页](https://github.com/xnervwang/AVMeta/releases) 下载对应的最新预编译版本使用。

**若使用预编译程序，可跳过此步骤**

1. 安装并配置 Golang + Git 开发环境， Golang 建议安装 1.13以上版本。
2. 执行命令：
    ```bash
   go get -u github.com/xnervwang/AVMeta
    ```
6. 至 `$GOPATH/bin` 目录下检查是否存在 `AVMeta` 可执行程序，并将 `$GOPATH/bin` 目录加入到环境变量中。

## 配置

在需要刮削的目录，执行命令 `AVMeta init` 生成 `config.yaml` 配置文件。

文件默认内容及解释如下：

```yaml
base:
  # 代理配置，格式为: socks5://127.0.0.1:1080 http://127.0.0.1:1080
  proxy: "socks5://127.0.0.1:1080"
media:
  # 媒体库配置，支持 nfo 和 vsmeta
  library: vsmeta
  # emby媒体库api访问地址，用于头像入库
  url: "http://127.0.0.1:8096"
  # emby媒体库api访问key
  api: ""
  # 腾讯云api id，用于面部识别裁图
  secretid: ""
  # 腾讯云api key，用于面部识别裁图
  secretkey: ""
path:
  # 刮削成功后存放的文件夹名称
  success: success
  # 刮削失败后存放的文件夹名称
  fail: fail
  # 刮削后所存放的路径
  # {actor} 演员中的第一个
  # {actors} 所有演员，以 "," 分隔
  # {number} 番号
  # {release} 发行日期
  # {year} 发行年份
  # {month} 发行月份
  # {studio} 厂商
  # {title} 电影名称
  # 比如下面的存放路径，番号为 "STARS-204",
  # 执行路径为 "/home/av"，最终保存的路径将会是
  # /home/av/success/SOD Create/2020/西野翔/STARS-204
  directory: '{studio}/{year}/{actor}/{number}'
  # 文件名中需要过滤的内容，以 "||" 分隔
  filter: -hd||hd-||[||]||【||】||asfur||~||-full||3xplanet||monv
site:
  # javbus免翻地址
  javbus: https://www.javbus.com/
  # javdb免翻地址
  javdb: https://javdb4.com/
```

## 使用

在使用之前请确保做了如下检查：

> 1. 已经成功安装 `AVMeta`
> 2. 已将 `$GOPATH/bin` 或 `AVMeta` 添加到环境变量
> 3. 确保需要刮削的视频文件均存放在程序执行目录下
> 4. 确保能够正常访问各类网站
> 5. 确保您所使用的账户对执行目录拥有读写权限
> 6. 最后，请确保在刮削目录下存在 **config.yaml** 配置文件，否则将使用默认配置

### 头像

本节仅针对 `emby` 媒体库用户，其余媒体库等待以后再说，若您所使用的不是 `emby` 媒体库，请跳过本节。

在入库头像之前，请您确保您的电脑能够正确访问 `emby` 媒体库，且您拥有一个 `api密钥`。

打开 `emby` 管理界面，并点击右上角 `管理` 按钮

![01.png](https://i.loli.net/2020/02/19/c2sT47Fw9XE8vMV.png)

点击左下角 `API 密钥` 按钮

![02.png](https://i.loli.net/2020/02/19/3qWFcxO4SujdeQg.png)

点击加号按钮创建 `API`

![03.png](https://i.loli.net/2020/02/19/v13Jh7QRBGzpuVT.png)

获取到 `API密钥` 后，请在配置文件中修改相应配置

#### 本地下载

本地下载头像，是将获取到的女优头像下载到本地，方便在后期无网络环境下也能入库。

若要下载女优头像，请在头像存放目录中执行命令:

```bash
AVMeta actress down
```

目前仅支持从 `javbus` 和 `javdb` 中获取女优头像。

默认命令将自动从两个网站下载所有女优头像，可通过添加 `--site javbus` `--site javdb` 参数来指定要下载的网站。

女优头像将保存在执行目录下的 `actress` 文件夹中，以 `女优名字.jpg` 的格式保存。

#### 本地入库

本地入库是方便本地存储有女优头像的朋友，在无需访问外网的情况下直接入库女优头像。

要执行本地入库，请先确保执行路径中存在 `actress` 文件夹，且文件夹中以 `女优名字.jpg` 格式存放有女优头像。

执行命令:

```bash
AVMeta actress put
```

入库时，程序会对女优名字进行搜索，若 `emby` 媒体库中存在此演员信息，且没有头像，则入库，反之不入库。

入库成功图片会移动到 `actress/sccess` 中。

### 刮削

刮削会根据从视频文件提取到的番号，自动搜索番号对应的元数据，并生成 *nfo* 或 *vsmeta* 元数据文件。

目前支持的元数据为：

- *nfo*: 对应 *emby*、*plex*、*kodi（未测试）*
- *vsmeta*: 对应群晖的 *DS Video* 系统

要对视频进行刮削，请将需要刮削的视频文件命名为正确的番号名称，并统一存放到一个目录中，然后执行命令：

```bash
AVMeta
```

#### NFO刮削

*nfo* 类型的元数据为通用元数据，无需特意指定媒体库程序。

在配置文件中修改 *Media* 下的 *Library* 为 *nfo*，即可生成 *nfo* 类型元数据。

生成后的元数据目录中将存放有： "视频文件"、"*.nfo"、"poster.jpg"、"fanart.jpg"

其中 "*.nfo" 为元数据信息文件，"poster.jpg" 为封面图片，"fanart.jpg" 为背景图片。

将生成后的元数据目录直接导入到对应的媒体库程序中，等待程序更新后即可查看。

#### 群晖刮削

*vsmeta* 类型的元数据为群晖 *DS Video* 使用元数据，仅支持群晖系统使用。

在配置文件中修改 *Media* 下的 *Library* 为 *vsmeta*，即可生成 *vsmeta* 类型元数据。

生成后的元数据目录中将存放有： "视频文件"、"*.vsmeta"

其中 "*.vsmeta" 为元数据信息文件，封面及背景图片也存储其中。

将生成的元数据目录直接导入到 *DS Video* 所设定的影片目录中，若无意外则等待更新后即可查看。

> PS: 若导入元数据后依然没有信息，请在 *DS Video* 设置中重建视频索引及视频信息，并在 *DS Video* 中将视频删除一次，再次导入等待更新。
> 这里需要注意，若在 *DS Video* 中删除视频，则对应视频文件及元数据也会一同删除，建议在本地保存一份再进行操作。

### 转换

若您原来使用的是 *nfo* 元数据文件，现今需要更换为 *DS Video*，那么可使用本程序提供的转换功能，自动将 *nfo* 元数据转换为 *DS Video* 所支持的群晖元数据文件。

在您的视频文件目录中，执行：

```bash
VSMeta nfo
```

程序将自动查询您当前目录下的所有存在 *nfo* 文件的目录，并自动将其进行转换。

请注意：

> 若 *.nfo* 同目录下存在 *fanart.jpg*、*poster.jpg* 文件，则会自动转换作为封面。
> 若不存在封面文件，则会通过自动下载 *.nfo* 中的封面信息进行转换。

## 鸣谢

特别感谢以下作者及所开发的程序，本项目参考过以下几位开发者代码及思想。

- [@yoshiko2](https://github.com/yoshiko2)，大部分设计思路及代码来源于 [AV_Data_Capture](https://github.com/yoshiko2/AV_Data_Capture)
- [@junerain123](https://github.com/junerain123)，人脸识别、图片裁剪来源于 [javsdt](https://github.com/junerain123/javsdt)
- [@soywiz](https://github.com/soywiz)，群晖支持部分来源于 [gist](https://gist.github.com/soywiz/2c10feb1231e70aca19a58aca9d6c16a)
