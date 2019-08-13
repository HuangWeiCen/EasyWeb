package GUI

import (
	"os"
	"bufio"
	"strings"
	"io"
	"io/ioutil"
)

// 读取一个端口的内容并保存到自身
func (this *PortRule) readConfig(constr string) error {
	startRead := false
	muxRule := MuxRule{Type: IsAFunc}
	read := ""
	cc := "" // 类型信息
	lastKuo := 0
	for i, s := range constr {
		if startRead {
			read += string(s)
		}
		// 处理括号
		// 判断一个端口的开始和结束
		if s == ']' {
			ind := strings.Index(read, cc)
			switch muxRule.Type {
			case IsAFunc:
				err := muxRule.analysis(read[8 : len(read)-1]) // 括号和标识(func:// )也被读到read中了，要把它去掉
				if err != nil {
					return err
				}
				break
			case IsAFileFunc:
				err := muxRule.readFile(strings.TrimSpace(read[ind : len(read)-1])) // 从类型标识之后读到右括号之前
				muxRule.Type = IsAFunc
				if err != nil {
					return err
				}
				break
			case IsAEasy:
				cs := strings.Split(strings.TrimSpace(read[ind+7:len(read)-1]), "=")
				if cs[0] == "EasyReceive" { // 上传需要的参数和别的不一样,要单独处理
					muxRule.Type = IsAReceive
					in := strings.Index(read, "(")
					if in != -1 {
						num := read[in+1 : len(read)-2]
						pat := read[ind+19 : in] // 去掉前面没用的部分,去掉easy://,去掉EasyReceive=
						r := Rule{}
						muxRule.analysis(`(path=` + pat + `,""="")`)
						r.analysis(num)
						muxRule.Rules = append(muxRule.Rules, r)
					}
					break
				}
				switch cs[0] {
				case "EasyFile":
					muxRule.Type = IsAFile
					break
				case "EasyFileAll":
					muxRule.Type = IsAFileAll
					break
				case "EsayDir":
					muxRule.Type = IsADir
					break
				}
				err := muxRule.analysis(`(path=` + cs[1] + `,""="")`) // 读取时获取path的值就好
				if err != nil {
					return err
				}
				break
			}

			this.Muxs = append(this.Muxs, muxRule)
			lastKuo = i
			startRead = false
			read = ""
			muxRule = MuxRule{Type: IsAFunc}
			cc = ""
			continue
		}
		if len(constr) >= i+2 && constr[i:i+2] == ":[" {
			muxRule.MuxPath = strings.TrimSpace(constr[lastKuo+1 : i])
			// 读取类型信息，去掉空格
			cc = strings.TrimSpace(constr[i+2 : i+9])
			for j := i; len(cc) < 7; j++ {
				cc += strings.TrimSpace(string(constr[j+9]))
			}
			switch cc {
			case "file://":
				muxRule.Type = IsAFileFunc
				break
			case "easy://":
				muxRule.Type = IsAEasy
				break
			}
			startRead = true

		}
	}
	return nil
}

// 读取规则文件并保存到自身
func (this *MuxRule) readFile(filePath string) error {
	err := IsRule(filePath)
	if err != IsRuleError {
		return err
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
				return nil
			}
			return err
		}
		if len(line) == 0 || line[:2] == "/*" || line == "rule" {
			continue
		}

		rule := Rule{}
		err = rule.analysis(line)
		if err != nil {
			return err
		}
		this.Rules = append(this.Rules, rule)

	}
}

//

// 解析一个路径并保存到自身
func (this *MuxRule) analysis(constr string) error {
	startRead := false
	rule := Rule{}
	read := ""
	for _, s := range constr {

		if startRead {
			read += string(s)
		}

		// 处理括号
		if s == '(' {
			startRead = true
			continue
		}
		if s == ')' {
			startRead = false
			rule.analysis(read[:len(read)-1]) // 去掉右括号
			this.Rules = append(this.Rules, rule)
			read = ""
			continue
		}
	}
	return nil
}

// 解析一条规则并保存到自身
func (this *Rule) analysis(str string) error {
	str = strings.TrimSpace(str)
	ss := strings.Split(str, ",")
	if len(ss) != 2 {
		return WrongFormatError
	}
	qs := strings.Split(ss[0], "=")
	as := strings.Split(ss[1], "=")
	if len(qs) != 2 || len(as) != 2 {
		return WrongFormatError
	}

	this.ComeKey = strings.TrimSpace(qs[0])
	this.ComeValue = strings.TrimSpace(qs[1])
	this.GoKey = strings.TrimSpace(as[0])
	govalue := strings.TrimSpace(as[1])
	if len(govalue) > 9 {
		if govalue[:9] != "return://" {
			this.GoValue = govalue
		} else {
			read, err := ioutil.ReadFile(govalue[9:])
			if err != nil {
				return err
			}
			this.GoValue = string(read)
		}
	} else {
		this.GoValue = govalue
	}
	return nil
}
