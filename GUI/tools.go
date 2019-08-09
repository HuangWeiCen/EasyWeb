package GUI

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"../Listener"
)

// 判断是规则文件还是配置文件并返回读取的内容，规则文件返回true
func IsRuleAndOpen(filePath string) error {
	if path.Ext(filePath) != ".easy" {
		return NotKnowFileError
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)

	for {

		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				return NotKnowFileError
			}
			return err
		}
		if line == "config" {
			return NotRuleError
		} else if line == "rule" {
			return IsRuleError
		}
	}
	return NotKnowFileError
}

// 读取文件并返回一个[]RortRule
func ReadConfigFile(filePath ...string) (prs []PortRule, err error) {
	for _, f := range filePath {
		err := IsRuleAndOpen(f)
		if err == NotRuleError {
			pr, err := readOneConfigFile(f)
			if err != nil {
				return nil, err
			}
			prs = append(prs, pr...)
		} else if err != nil {
			return nil, err
		}
	}
	return prs, err
}

// 读取一个文件的内容并返回一个[]RortRule
func readOneConfigFile(file string) (prs []PortRule, err error) {
	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	pr := PortRule{}
	var startRead bool
	var read string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				return prs, nil
			}
			return prs, err
		}

		// 跳过注释
		if line[:2] == "/*" {
			continue
		}

		// 读取内容,先判断要不要读
		if startRead {
			read += line
		}

		// 判断一个端口的开始和结束
		if line[len(line)-1] == '}' {
			read += line[:len(line)-1]
			err = pr.readConfig(read[:len(read)-1]) // 右括号也被读到read中了，要把它去掉
			if err != nil {
				return nil, err
			}
			prs = append(prs, pr)
			startRead = false
			read = ""
			pr = PortRule{}
			continue
		}
		if i := strings.Index(line, ":{"); i != -1 {
			pr.PortPath = ":" + line[:i]
			startRead = true
			read += line[i:]
			continue
		}
	}

	return prs, nil
}

// 根据文件创建EasyListen对象
func NewEasyListenToPortRuleList(ip string, portRules ...PortRule) *Listener.EasyListen {
	listen := Listener.NewEasyListener(ip)
	for _, p := range portRules { // 设置端口
		listen.AddPort(p.PortPath)
		for _, m := range p.Muxs { // 设置端口对应路由
			switch m.Type {
			case IsAReceive:
				max, path, key := "", "", ""
				for _, r := range m.Rules {
					switch r.ComeKey {
					case "path":
						path = r.ComeValue.(string)
						break
					case "max":
						max = r.ComeValue.(string)
						key = r.GoValue.(string)
					case "key":
						key = r.ComeValue.(string)
						max = r.GoValue.(string)
					}
				}
				i6, err := strconv.ParseInt(max, 10, 64)
				if err != nil {
					break
				}
				listen.AddPostReceiveMux(i6, p.PortPath, key, m.MuxPath, path)
				break
			case IsADir:
				listen.AddDirMux(p.PortPath, m.MuxPath, m.Rules[0].ComeValue.(string))
				break
			case IsAFileAll:
				listen.AddDirMux(p.PortPath, m.MuxPath, m.Rules[0].ComeValue.(string))
				break
			case IsAFile:
				listen.AddDirMux(p.PortPath, m.MuxPath, m.Rules[0].ComeValue.(string))
				break
			case IsAFunc:
				listen.AddEasyFuncMux(p.PortPath, m.MuxPath, &postListener{MuxRule{Rules: m.Rules}})
				break
			}
		}
	}
	return listen
}

// 重写方法
func (this *postListener) EasyHttpListen(values url.Values, responseType Listener.ResponseType) []byte {
	mp := make(map[string]interface{})
	if responseType == Listener.EPOST {
		for _, r := range this.Rules {
			if values.Get(r.ComeKey) == r.ComeValue {
				mp[r.GoKey] = r.GoValue
			}
		}
	}
	bts, err := json.Marshal(mp)
	if err != nil {
		fmt.Errorf("出现异常: %s", err)
	}
	return bts
}
