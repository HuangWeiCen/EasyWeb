package GUI

import (
	"path"
	"os"
	"bufio"
	"strings"
	"io"
	"io/ioutil"
)

// 判断是规则文件还是配置文件并返回读取的内容，规则文件返回true
func IsRuleAndOpen(filePath string) (*os.File, error) {
	if path.Ext(filePath) != ".easy" {
		return nil, NotKnowFileError
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				return nil, NotKnowFileError
			}
			return nil, err
		}
		if line == "config" {
			return f, NotRuleError
		} else if line == "rule" {
			return f, IsRuleError
		}
	}
	return nil, NotKnowFileError
}

// 读取文件并保存到自己内
func ReadConfigFile(filePath ...string) (prs []PortRule, err error) {
	for _, f := range filePath {
		file, err := IsRuleAndOpen(f)
		if err == NotRuleError {
			pr, err := readOneConfigFile(file)
			if err != nil {
				return nil, err
			}
			prs = append(prs, pr...)
		} else if err != nil {
			return nil, err
		}
	}
	return
}

// 读取一个文件的内容并返回一个[]RortRule
func readOneConfigFile(file *os.File) (prs []PortRule, err error) {
	buf := bufio.NewReader(file)
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
			return
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
	return
}

// 读取一个端口的内容并保存到自身
func (this *PortRule) readConfig(constr string) error {
	startRead := false
	muxRule := MuxRule{Type: IsAFunc}
	read := ""
	cc := ""
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
				err := muxRule.analysis(read[:len(read)-1]) // 右括号也被读到read中了，要把它去掉
				if err != nil {
					return err
				}
				break
			case IsAFileFunc:
				err := muxRule.readFile(strings.TrimSpace(read[ind : len(read)-1])) // 从类型标识之后读到右括号之前
				if err != nil {
					return err
				}
				break
			case IsAEasy:
				cs := strings.Split(strings.TrimSpace(read[ind:len(read)-1]), "=")
				switch cs[0] {
				case "EasyFile":
					muxRule.Type = IsAFile
					break
				case "EasyFileAll":
					muxRule.Type = IsAFileAll
					break
				case "EasyDir":
					muxRule.Type = IsADir
					break
				case "EasyReceive":
					muxRule.Type = IsAReceive
					break
				}
				err := muxRule.analysis(`(path=` + cs[1] + `,""="")`) // 读取时获取path的值就好
				if err != nil {
					return err
				}
			}

			this.Muxs = append(this.Muxs, muxRule)
			startRead = false
			read = ""
			muxRule = MuxRule{Type: IsAFunc}
			cc = ""
			continue
		}
		if constr[i-1:i+1] == ":[" {
			this.PortPath = constr[:i-1]
			// 读取类型信息，去掉空格
			cc = strings.TrimSpace(constr[i+1 : i+8])
			for j := i; len(cc) < 7; j++ {
				cc += strings.TrimSpace(string(constr))
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
		if len(line) == 0 {
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
	ss := strings.Split(str, "=")
	if len(ss) != 2 {
		return WrongFormatError
	}
	qs := strings.Split(ss[0], ",")
	as := strings.Split(ss[1], ",")
	if len(qs) != 2 || len(as) != 2 {
		return WrongFormatError
	}

	this.ComeKey = strings.TrimSpace(qs[0])
	this.ComeValue = strings.TrimSpace(qs[1])
	this.GoKey = strings.TrimSpace(as[0])
	govalue := strings.TrimSpace(as[1])
	if govalue[:9] != "return://" {
		this.GoValue = govalue
	} else {
		read, err := ioutil.ReadFile(govalue[9:])
		if err != nil {
			return err
		}
		this.GoValue = read
	}
	return nil
}
