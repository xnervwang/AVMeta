package cmd

import (
	"encoding/xml"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xnervwang/AVMeta/pkg/logs"
	"github.com/xnervwang/AVMeta/pkg/media"
	"github.com/xnervwang/AVMeta/pkg/util"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	// 定义扩展列表
	videoExts = map[string]string{
		".avi":  ".avi",
		".flv":  ".flv",
		".mkv":  ".mkv",
		".mov":  ".mov",
		".mp4":  ".mp4",
		".rmvb": ".rmvb",
		".ts":   ".ts",
		".wmv":  ".wmv",
	}
)

// NfoFile nfo文件列表结构
type NfoFile struct {
	Path  string
	Video string
	Dir   string
}

// nfo命令
func (e *Executor) initNfo() {
	nfoCmd := &cobra.Command{
		Use: "nfo",
		Long: `
自动将运行目录下所有nfo文件转换为VSMeta文件`,
		Example: `  AVMeta nfo`,
		Run:     e.nfoRunFunc,
	}

	e.rootCmd.AddCommand(nfoCmd)
}

// 转换执行命令
func (e *Executor) nfoRunFunc(cmd *cobra.Command, args []string) {
	// 初始化日志
	logs.Log("")

	// 获取当前执行路径
	curDir := util.GetRunPath()

	// 文件列表
	var nfos []NfoFile

	// 列当前目录
	nfos, err := e.walk(curDir, nfos)
	// 检测错误
	logs.FatalError(err)

	// 获取总量
	count := len(nfos)
	// 输出总量
	logs.Info("共探索到 %d 个 nfo 文件, 开始转换...\n\n", count)

	// 初始化进程
	wg := util.NewWaitGroup(2)

	// 循环nfo文件列表
	for _, nfo := range nfos {
		// 计数加
		wg.AddDelta()
		// 转换进程
		go e.nfoProcess(nfo, wg)
	}

	// 等待结束
	wg.Wait()
}

// 转换进程
func (e *Executor) nfoProcess(nfo NfoFile, wg *util.WaitGroup) {
	// 读取文件
	b, err := util.ReadFile(nfo.Path)
	// 检查
	if err != nil {
		// 输出错误
		logs.Error("文件: [%s] 打开失败, 错误原因: %s\n", path.Base(nfo.Path), err)

		// 进程
		wg.Done()

		return
	}

	// 媒体对象
	var m media.Media

	// 转换
	err = xml.Unmarshal(b, &m)
	// 检查错误
	if err != nil {
		// 输出错误
		logs.Error("文件: [%s] 打开失败, 错误原因: %s\n", path.Base(nfo.Path), err)

		// 进程
		wg.Done()

		return
	}

	// 实例化vsmeta
	vs := media.NewVSMeta()
	// fanart
	if !util.Exists(nfo.Dir+"/fanart.jpg") && m.FanArt != "" {
		err = util.SavePhoto(m.FanArt, fmt.Sprintf("%s/fanart.jpg", nfo.Dir), "", !strings.EqualFold(strings.ToLower(path.Ext(m.FanArt)), ".jpg"))
		if err != nil {
			// 输出警告
			logs.Warning("文件: [%s] 封面下载失败, 错误原因: %s\n", path.Base(nfo.Path), err)
		}
	}
	// poster
	if !util.Exists(nfo.Dir+"/poster.jpg") && m.Poster != "" {
		err = util.SavePhoto(m.Poster, fmt.Sprintf("%s/poster.jpg", nfo.Dir), "", !strings.EqualFold(strings.ToLower(path.Ext(m.Poster)), ".jpg"))
		if err != nil {
			// 输出错误
			logs.Warning("文件: [%s] 封面下载失败, 错误原因: %s\n", path.Base(nfo.Path), err)
		}
	}

	m.FanArt = fmt.Sprintf("%s/fanart.jpg", nfo.Dir)
	m.Poster = fmt.Sprintf("%s/poster.jpg", nfo.Dir)

	// 解析为 vsmeta
	bs := vs.Convert(&m)

	// 获取视频后缀
	ext := path.Ext(nfo.Video)

	// 写入vsmeta
	err = util.WriteFile(fmt.Sprintf("%s/%s%s.vsmeta", nfo.Dir, m.Number, ext), bs)
	// 检查
	if err != nil {
		// 输出错误
		logs.Error("文件: [%s] 转换失败, 错误原因: %s\n", path.Base(nfo.Path), err)

		// 进程
		wg.Done()

		return
	}

	// 输出正确
	logs.Info("文件: [%s/%s] 转换成功, 路径: %s\n", path.Base(nfo.Path), m.Number, nfo.Dir)

	// 进程
	wg.Done()
}

// 列目录
func (e *Executor) walk(dirPath string, nfoFiles []NfoFile) ([]NfoFile, error) {
	// 读取目录
	r, err := ioutil.ReadDir(dirPath)
	// 检查错误
	if err != nil {
		return nil, err
	}

	// 循环列表
	for _, f := range r {
		if f.IsDir() {
			fullDir := dirPath + "/" + f.Name()
			nfoFiles, err = e.walk(fullDir, nfoFiles)
		} else {
			// 获取后缀并转换为小写
			ext := strings.ToLower(path.Ext(f.Name()))

			// 验证是否存在于后缀扩展名中
			if _, ok := videoExts[ext]; ok {
				// 遍历目录
				err := filepath.Walk(dirPath, func(filePath string, fi os.FileInfo, err error) error {
					// 错误
					if fi == nil {
						return err
					}

					if !fi.IsDir() {
						// 是否nfo文件
						if strings.ToLower(filepath.Ext(fi.Name())) == ".nfo" {
							// 初始化nfo
							nfo := NfoFile{
								Path:  dirPath + "/" + fi.Name(),
								Video: dirPath + "/" + f.Name(),
								Dir:   dirPath,
							}

							// 加入列表
							nfoFiles = append(nfoFiles, nfo)
						}
					}

					return nil
				})

				if err != nil {
					return nfoFiles, err
				}
			}
		}
	}

	return nfoFiles, nil
}
