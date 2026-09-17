# m3u8-downloader

golang 多线程下载直播流m3u8格式的视屏，跨平台。 你只需指定必要的 flag (`u`、`o`、`n`、`ht`) 来运行, 工具就会自动帮你解析 M3U8 文件，并将 TS 片段下载下来合并成一个文件。


## 功能介绍

1. 下载和解析 M3U8
2. 下载 TS 失败重试 （加密的同步解密)
3. 合并 TS 片段
4. 批量下载（CSV 文件指定多个地址）
5. 独立运行 exe 时，若无任何地址源，会交互式提示输入 m3u8 地址，便于普通用户使用

> 可以下载岛国小电影  
> 可以下载岛国小电影  
> 可以下载岛国小电影    
> 重要的事情说三遍......

## 效果展示
![demo](./demo.gif)

## 项目结构

```
.
├── m3u8-downloader.go       # 核心下载库（包功能，供导入调用，不含 CLI 逻辑）
├── cmd/m3u8-downloader/
│   └── main.go              # 命令行入口（仅独立 exe 使用）
├── build-release.sh         # 多平台打包脚本
└── m3u8-downloader_test.go  # 测试
```

包功能与 exe 功能各自独立：给包调用方修改功能不影响 exe 行为，修改 CLI 也不必动库。

## 参数说明：

```
- u M3U8 地址
- f 包含多个m3u8地址的 csv(url,filename) 文件路径
- o 自定义文件名, 默认 movie
- n 下载协程并发数，默认 16
- ht 设置getHost的方式（共两种 apiv1 和 apiv2）, 默认 apiv1
- c 自定义请求cookie, 默认空
- s 是否允许不安全的请求, 默认 0
- sp 文件保存路径, 默认为当前路径
```

默认情况只需要传`u`参数,其他参数保持默认即可。 部分链接可能限制请求频率，可根据实际情况调整 `n` 参数的值。

不带任何参数直接运行 exe（双击运行）时，会提示输入 m3u8 地址，回车即可开始下载；直接回车则退出。

## 下载

已经编译好的平台有： [点击下载](https://github.com/llychao/m3u8-downloader/releases)

- windows/amd64
- linux/amd64
- darwin/amd64

## 用法

### 源码方式

```bash
自己编译：go build -o m3u8-downloader ./cmd/m3u8-downloader
简洁使用：./m3u8-downloader  -u=http://example.com/index.m3u8
批量使用：./m3u8-downloader  -f=videos.csv
完整使用：./m3u8-downloader  -u=http://example.com/index.m3u8 -o=example -n=16 -ht=apiv1 -c="key1=v1; key2=v2"
```

### 作为库导入

示例代码：

```go
package main

import (
    "log"

    m3u8downloader "github.com/Linchpin-L/m3u8-downloader"
)

func main() {
    // 简洁方式：使用默认参数（16 线程等），merge 为 true 时文件名会自动携带 .ts 后缀
    if err := m3u8downloader.Download("https://example.com/index.m3u8", "example", true); err != nil {
        log.Fatal(err)
    }

    // 不合并分片：filename 作为文件夹名，ts 分片保留在该文件夹内
    if err := m3u8downloader.Download("https://example.com/index.m3u8", "example", false); err != nil {
        log.Fatal(err)
    }

    // 完整方式：自定义线程数、host 方式、cookie、保存路径等（最后一个参数为是否合并）
    m3u8downloader.DownloadSingleVideo("https://example.com/index.m3u8", 16, "", "example", "", 0, "", 0, true)

    // 直接下载非 m3u8 资源（例如 mp4）
    m3u8downloader.DownloadDirect("https://example.com/video.mp4", "video", "", 0, "")
}
```

### 二进制方式:

Linux 和 MacOS 和 Windows PowerShell

```
简洁使用：
./m3u8-downloader-v1.0.0-linux-amd64 -u=http://example.com/index.m3u8
./m3u8-downloader-v1.0.0-darwin-amd64 -u=http://example.com/index.m3u8 
.\m3u8-downloader-v1.0.0-windows-amd64.exe -u=http://example.com/index.m3u8

批量使用：
./m3u8-downloader-v1.0.0-linux-amd64 -f=videos.csv

完整使用：
./m3u8-downloader-v1.0.0-linux-amd64 -u=http://example.com/index.m3u8 -o=example -n=16 -ht=apiv1 -c="key1=v1; key2=v2"
./m3u8-downloader-v1.0.0-darwin-amd64 -u=http://example.com/index.m3u8 -o=example -n=16 -ht=apiv1 -c="key1=v1; key2=v2"
.\m3u8-downloader-v1.0.0-windows-amd64.exe -u=http://example.com/index.m3u8 -o=example -n=16 -ht=apiv1 -c="key1=v1; key2=v2"

交互使用（不带任何参数运行）：
.\m3u8-downloader-v1.0.0-windows-amd64.exe
未指定下载地址源，请输入 m3u8 视频地址(http(s):// 开头，直接回车退出): https://example.com/index.m3u8
```

## TODO

- [ ] 合并结果兼容性问题：当前使用 `copy /b`（Windows）或 `cat`（Unix）对 TS 分片做纯字节拼接，不重写分片内的 PTS/DTS 时间戳。分片各自携带独立时间轴，拼接后时间戳不连续——PotPlayer、VLC 等兼容性强的播放器可正常播放，但 Windows 自带"电影和电视"等基于 Media Foundation 的播放器通常只能播放开头几秒。计划改用 ffmpeg remux（`ffmpeg -i merge.ts -c copy out.mp4`）或自行重写时间轴解决。
- [ ] 分片静默缺失：单个 TS 分片重试多次仍失败时会被跳过，合并不报错，最终文件中间会缺一段。计划增加缺失分片检测，合并前给出警告。

## 问题说明

1.在Linux或者mac平台，如果显示无运行权限，请用chmod 命令进行添加权限
```bash
 # Linux amd64平台
 chmod 0755 m3u8-downloader-v1.0.0-linux-amd64
 # Mac darwin amd64平台
 chmod 0755 m3u8-downloader-v1.0.0-darwin-amd64
 ```
2.下载失败的情况,请设置 -ht="apiv1" 或者 -ht="apiv2" （默认为apiv1）
```golang
func get_host(Url string, ht string) string {
    u, err := url.Parse(Url)
    var host string
    checkErr(err)
    switch ht {
    case "apiv1":
        host = u.Scheme + "://" + u.Host + path.Dir(u.Path)
    case "apiv2":
        host = u.Scheme + "://" + u.Host
    }
    return host
}
```
