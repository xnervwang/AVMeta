package media

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/xnervwang/AVMeta/pkg/scraper"
)

// Media Nfo信息结构，
// 用以存储 nfo 文件所需各项信息。
type Media struct {
	XMLName   xml.Name `xml:"movie"`
	Title     Inner    `xml:"title"`
	SortTitle string   `xml:"sorttitle"`
	Number    string   `xml:"num"`
	Studio    Inner    `xml:"studio"`
	Maker     Inner    `xml:"maker"`
	Director  Inner    `xml:"director"`
	Release   string   `xml:"release"`
	Premiered string   `xml:"premiered"`
	Year      string   `xml:"year"`
	Plot      Inner    `xml:"plot"`
	Outline   Inner    `xml:"outline"`
	RunTime   string   `xml:"runtime"`
	Mpaa      string   `xml:"mpaa"`
	Country   string   `xml:"country"`
	Poster    string   `xml:"poster"`
	Thumb     string   `xml:"thumb"`
	FanArt    string   `xml:"fanart"`
	Actor     []Actor  `xml:"actor"`
	Tag       []Inner  `xml:"tag"`
	Genre     []Inner  `xml:"genre"`
	Set       string   `xml:"set"`
	Label     string   `xml:"label"`
	Cover     string   `xml:"cover"`
	WebSite   string   `xml:"website"`
	Month     string   `xml:"-"`
	DirPath   string   `xml:"-"`
	Source    string   `xml:"-"`
}

// Inner 文字数据，为了避免某些内容被转义。
type Inner struct {
	Inner string `xml:",innerxml"`
}

// Actor 演员信息，保存演员姓名及头像地址。
type Actor struct {
	Name  string `xml:"name"`
	Thumb string `xml:"thumb"`
}

// ParseMedia 将刮削对象解析为 Media 结构体，
// 解析错误时返回空对象及错误信息。
//
// s IScraper刮削接口，传入刮削对象
// site 字符串参数，传入刮削网站
func ParseMedia(s scraper.IScraper, site string) (*Media, error) {
	// 定义一个nfo对象
	var m Media

	// 设定刮削网站
	m.Source = site

	// 检查刮削对象
	if s == nil {
		return nil, fmt.Errorf("scraper no data")
	}

	// 定义演员列表
	var actors []Actor

	// 获取演员并循环
	for name, thumb := range s.GetActors() {
		// 加入列表
		actors = append(actors, Actor{
			Name:  name,
			Thumb: thumb,
		})
	}

	// 短标题
	m.SortTitle = strings.TrimSpace(s.GetNumber())
	// 番号
	m.Number = m.SortTitle
	// 厂商
	m.Studio = Inner{Inner: strings.TrimSpace(s.GetStudio())}
	// 厂商
	m.Maker = m.Studio
	// 导演
	m.Director = Inner{Inner: strings.TrimSpace(s.GetDirector())}
	// 发行时间
	m.Release = strings.TrimSpace(strings.ReplaceAll(s.GetRelease(), "/", "-"))
	// 发行时间
	m.Premiered = m.Release
	// 设置年份
	m.Year = strings.TrimSpace(GetYear(m.Release))
	// 简介
	m.Plot = Inner{Inner: s.GetIntro()}
	// 简介
	m.Outline = m.Plot
	// 时长
	m.RunTime = strings.TrimSpace(s.GetRuntime())
	// 分级
	m.Mpaa = "XXX"
	// 国家
	m.Country = "JP"
	// 演员
	m.Actor = actors
	// 标签
	tags := s.GetTags()
	// 循环标签
	for _, tag := range tags {
		m.Tag = append(m.Tag, Inner{Inner: tag})
	}
	// 类型
	m.Genre = m.Tag
	// 系列
	m.Set = strings.TrimSpace(s.GetSeries())
	// 图片
	m.Cover = strings.TrimSpace(s.GetCover())
	// 地址
	m.WebSite = strings.TrimSpace(s.GetURI())
	// 设置月份
	m.Month = strings.TrimSpace(GetMonth(m.Release))

	// 获取标题
	title := strings.TrimSpace(s.GetTitle())
	// 替换原有番号
	title = strings.TrimSpace(strings.ReplaceAll(title, m.Number, ""))
	// 重新增加番号
	title = fmt.Sprintf("%s %s", m.Number, title)
	// 设置标题
	m.Title = Inner{Inner: title}

	return &m, nil
}

// GetYear 通过获取到的发行日期获取年份信息。
//
// date 字符串参数，传入发行日期。
func GetYear(date string) string {
	// 年份搜索正则
	re := regexp.MustCompile(`\d{4}`)

	return re.FindString(date)
}

// GetMonth 通过获取到的发行日期获取月份信息。
//
// date 字符串参数，传入发行日期。
func GetMonth(date string) string {
	// 月份搜索正则
	re := regexp.MustCompile(`\d{4}-([\d]{2})-\d{2}`)
	// 查找
	month := re.FindStringSubmatch(date)
	// 如果找到
	if len(month) > 0 {
		return month[1]
	}

	return ""
}

// ConvertMap 将部分内容转换为 map 对象，
// 该方法主要用于路径配置中的数据转换。
func (m *Media) ConvertMap() map[string]string {
	// 定义map
	replaceMap := make(map[string]string)
	// 演员数组
	var actors []string
	// 循环
	for _, actor := range m.Actor {
		// 加入数组
		actors = append(actors, actor.Name)
	}
	// 是否有演员
	if len(actors) > 0 {
		replaceMap["{actor}"] = actors[0]
	} else {
		replaceMap["{actor}"] = "未知演员"
	}
	// 替换演员名单
	replaceMap["{actors}"] = strings.Join(actors, ",")
	// 替换番号
	replaceMap["{number}"] = m.Number
	// 替换发行时间
	replaceMap["{release}"] = m.Release
	// 替换年份
	replaceMap["{year}"] = m.Year
	// 替换月份
	replaceMap["{month}"] = m.Month
	// 替换影片公司
	replaceMap["{studio}"] = m.Studio.Inner
	// 替换影片名称
	replaceMap["{title}"] = m.Title.Inner

	return replaceMap
}
