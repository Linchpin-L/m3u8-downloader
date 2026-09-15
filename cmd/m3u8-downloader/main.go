// @author:llychao<lychao_vip@163.com>
// @功能:m3u8-downloader 命令行入口（仅独立 exe 使用，包功能见根目录 m3u8downloader）
package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	m3u8downloader "github.com/Linchpin-L/m3u8-downloader"
)

// InputEntry 用于保存 -f 中的 CSV 行
type InputEntry struct {
	URL      string
	Filename string
}

var (
	// 命令行参数
	urlFlag = flag.String("u", "", "m3u8下载地址(http(s)://url/xx/xx/index.m3u8)")
	nFlag   = flag.Int("n", 16, "下载线程数(max goroutines num)")
	htFlag  = flag.String("ht", "", "设置getHost的方式(apiv1: `http(s):// + url.Host + filepath.Dir(url.Path)`; apiv2: `http(s)://+ u.Host`")
	oFlag   = flag.String("o", "movie", "自定义文件名(默认为movie)")
	cFlag   = flag.String("c", "", "自定义请求 cookie")
	sFlag   = flag.Int("s", 0, "是否允许不安全的请求(默认为0)")
	spFlag  = flag.String("sp", "", "文件保存路径(默认为当前路径)")
	fFlag   = flag.String("f", "", "包含多个m3u8地址的 csv(url,filename) 文件路径")
)

func main() {
	Run()
}

func Run() {
	msgTpl := "[功能]:多线程下载直播流 m3u8 视频（ts + 合并）\n[提醒]:如果下载失败，请使用 -ht=apiv2 \n[提醒]:如果下载失败，m3u8 地址可能存在嵌套\n[提醒]:如果进度条中途下载失败，可重复执行"
	fmt.Println(msgTpl)
	fmt.Println("[开发]:linchpin1029@qq.com, 版本：1.0.0")
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 解析命令行参数
	flag.Parse()
	m3u8Url := *urlFlag
	filePath := *fFlag
	maxGoroutines := *nFlag
	hostType := *htFlag
	movieDir := *oFlag
	cookie := *cFlag
	insecure := *sFlag
	savePath := *spFlag

	var entries []InputEntry
	if filePath != "" {
		// 从 CSV 文件读取 URL 和 Filename
		f, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Error] 无法读取文件 %s: %v\n", filePath, err)
			os.Exit(1)
		}
		r := csv.NewReader(f)
		records, err := r.ReadAll()
		f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Error] 无法解析 CSV %s: %v\n", filePath, err)
			os.Exit(1)
		}
		for _, rec := range records {
			if len(rec) < 2 {
				continue
			}
			u := strings.TrimSpace(rec[0])
			fn := strings.TrimSpace(rec[1])
			// skip header if present
			if strings.ToLower(u) == "url" && strings.ToLower(fn) == "filename" {
				continue
			}
			if u == "" {
				continue
			}
			entries = append(entries, InputEntry{URL: u, Filename: fn})
		}
	}
	if m3u8Url != "" && filePath == "" {
		// 单个 URL (来自 -u)，使用 -o 作为文件名
		entries = append(entries, InputEntry{URL: m3u8Url, Filename: movieDir})
	}

	if len(entries) == 0 {
		// 独立运行 exe 时没有任何地址源，交互式提示输入
		flag.Usage()
		u := promptM3u8URL()
		if u == "" {
			return
		}
		entries = append(entries, InputEntry{URL: u, Filename: movieDir})
	}

	// 循环下载每个条目（支持 m3u8 和 直接文件）
	for i, entry := range entries {
		if len(entries) > 1 {
			fmt.Printf("\n========== 开始下载第 %d/%d 个视频 ==========_\n", i+1, len(entries))
		}
		// 判断是否 m3u8（根据 URL 路径后缀）
		isM3u8 := false
		if u, err := url.Parse(entry.URL); err == nil {
			if strings.HasSuffix(strings.ToLower(u.Path), ".m3u8") {
				isM3u8 = true
			}
		} else if strings.HasSuffix(strings.ToLower(entry.URL), ".m3u8") {
			isM3u8 = true
		}

		if isM3u8 {
			// 确保传入的目录名不带扩展，最终会使用 .mp4
			nameOnly := strings.TrimSuffix(entry.Filename, filepath.Ext(entry.Filename))
			m3u8downloader.DownloadSingleVideo(entry.URL, maxGoroutines, hostType, nameOnly, cookie, insecure, savePath, i)
		} else {
			// 直接下载（例如 mp4）并保持原后缀（如果 filename 未包含后缀，则补齐 URL 的后缀）
			if err := m3u8downloader.DownloadDirect(entry.URL, entry.Filename, cookie, insecure, savePath); err != nil {
				fmt.Fprintf(os.Stderr, "[Error] 下载文件失败 %s: %v\n", entry.URL, err)
			}
		}
	}

	waitExit()
}

// waitExit 阻塞等待用户按回车，避免双击运行时窗口在输出可见前关闭
func waitExit() {
	fmt.Print("\n按回车键退出...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// promptM3u8URL 独立 exe 无任何地址源时，交互式提示用户输入 m3u8 地址
// 返回空串表示用户放弃输入
func promptM3u8URL() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n未指定下载地址源，请输入 m3u8 视频地址(http(s):// 开头，直接回车退出): ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return ""
		}
		line = strings.TrimSpace(line)
		if line == "" {
			return ""
		}
		if !strings.HasPrefix(strings.ToLower(line), "http") {
			fmt.Println("[提醒] 地址需以 http(s):// 开头，请重新输入")
			continue
		}
		return line
	}
}
