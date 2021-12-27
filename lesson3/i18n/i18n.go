package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 本地化包
var bundle *i18n.Bundle

func init() {
	// 设置默认语言
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	//os.Args是用来获取命令行执行参数分片的，
	//LookPath去环境变量中找这个可执行文件的绝对路径，或相对于当前目录的路径
	file, _ := exec.LookPath(os.Args[0])
	//根据传入的路径计算出绝对路径，如果传入的为相对路径，那么它会把当前路径拼接上
	path, _ := filepath.Abs(file)
	//path是一个包含可执行文件在内的完整路径，我们只需要精确到目录即可
	//搜索最后一个目录分隔符的位置（下标），然后通过以下代码将路径中下标后面的字符串切割掉
	index := strings.LastIndex(path, string(os.PathSeparator))
	//完成了目录的获取，接下来再拼接上我们实际的配置文件就可以了
	path = path[:index]

	// 加载公共信息语言包文件
	bundle.MustLoadMessageFile(path + "/tomls/nft.en.toml")
	bundle.MustLoadMessageFile(path + "/tomls/nft.zh-cn.toml")
}

/* 本地化语言
param:
	lang 语言 en|zh-cn
	messageID 语言文件中的 messageID
	errCode: string 错误码
	templateData: i18n 文件中需要变量替换的内容
	pluralCount: 传入 int 或 int64 类型数据 ， 根据数字判断是否返回复数格式 msg
*/
func MustLocalize(lang, accept, messageID string, templateData interface{}, pluralCount interface{}) string {
	localizer := i18n.NewLocalizer(bundle, lang, accept)
	return localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  pluralCount,
	})
}
