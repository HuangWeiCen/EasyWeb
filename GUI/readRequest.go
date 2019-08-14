package GUI

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"os"
	"strings"
)

/*
请使用json格式发送数据
json报文一:
	Ip:设置的ip地址(string)
json报文二:--重复发一个路由地址可能会冲掉前面的地址(注意可能)
	Port:下方路由地址对应的端口号(string)
	Mux:设置的路由地址(string)
	Type:路由地址的类型(string)--映射就是能从网络上看到并下载
		--普通:Normal ; 上传:Upload ; 文件映射:File ; 映射文件夹内容(不含文件夹):FileAll ; 映射文件夹:Dir
	规则数据:规则数据的键随便起,两个同名的为一组,键的结尾分别加上1和2.
		--1结尾的为用户可能发过来的键值对,2为对应的以1结尾的键值对出现时服务器要返回的数据.
		--规则数据的值写成键=值的形式 比如{"such1":"server=hello","such2":"client=hi"}
		--请成对出现,如果数据没发全可能会产生任何异常情况(我做了预防,但肯定有没想到的地方)
	Path:路径信息(string) 普通模式时会读取该地址下的配置文件,其他模式下会将此地址设置成对应地址
	Max:上传最大值(int) 上传模式时设置的最大上传大小
	Upkey:上传时使用的key--暂时只支持post表单传输
json报文三:
	"Start":"ok"
	固定写法,用于启动监听
备注:不可以混着写,请注意报文二该传的数据要传全,传不全此条数据可能会生效一半导致出现各种问题


返回报文:
	Set: "T" or "F" T成功,F失败
	Reason: 错误原因,就那么固定几条,看那几个SetErrorString就好
 */

// 解析and监听
func JieXiAndJianTing(port, mux string) {
	if mux[0] != '/' {
		mux = "/" + mux
	}
	PortRules = []PortRule{}
	http.HandleFunc(mux, jieXi1)
	http.ListenAndServe(":"+port, nil)
}

// 解析控制界面发送的数据
// 早晚得重构,一个方法写的太长了,还很难分成几段
func jieXi1(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var cip string                       // IP地址
	var cport string                     // 端口号
	var cmux string                      // 路由地址
	var cpath string                     // 路径地址
	var ctype string                     // 路由类型
	var cupkey string                    // 上传的键
	var cmax string                      // 最大上传
	var rules map[string]Rule            // 规则map,key为用户发送的键
	var rebt []byte                      // 返回的时候必须以这个格式返回,先建好变量
	returnmap := make(map[string]string) // 返回成功与否和失败原因

	errF := func(errstr SetErrorString, err error) {
		if err != nil {
			fmt.Printf("出现异常 : %s", err)
		}
		if errstr == SuccessErrStr {
			returnmap["Set"] = "T"
		} else {
			returnmap["Set"] = "F"
		}
		returnmap["Reason"] = string(errstr)
		rebt, err = json.Marshal(returnmap)
		if err != nil {
			rebt, err = json.Marshal(returnmap)
			for {
				if err != nil {
					fmt.Printf("出现异常 : %s", err)
					returnmap["Set"] = "F"
					returnmap["Reason"] = string(SetErrStr)
					rebt, err = json.Marshal(returnmap)
				} else {
					w.Write(rebt)
					return
				}
			}
		} else {
			w.Write(rebt)
			return
		}
	} // 异常处理
	havePath := func(m map[string]interface{}) bool {
		if m["Path"] != nil {
			_, err := os.Stat(m["Path"].(string))
			if err != nil {
				errF(RuleErrStr, err)
				return false
			}
			cpath = m["Path"].(string)
			return true
		} else {
			errF(RuleMissingErrStr, nil)
			return false
		}
	} // 判断填没填地址,地址有没有填错

	if r.Method == "POST" {
		read, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			fmt.Errorf("出现异常：%s", err)
			return
		}
		m := make(map[string]interface{})
		err = json.Unmarshal(read, &m)
		switch { // 这算不算go的语法糖
		case m["Ip"] != nil:
			cip = m["Ip"].(string)
		case m["Port"] != nil:
			cport = m["Port"].(string)
			fallthrough
		case m["Mux"] != nil:
			cmux = m["Mux"].(string)
			if cport == "" {
				errF(RuleErrStr, nil)
			}
			fallthrough
			// 稀奇古怪的go语言,下面的这段代码会报:"interface conversion: interface {} is nil, not string"
			// 蜜汁程序,先判定为非空,再一本正经的告诉我这是空的无法转换类型
			/*case m["Path"] != nil:
				cpath = m["Path"].(string)
				fallthrough*/
		case m["Type"] != nil:
			if m["Type"] != nil {//这个语法糖跟搞笑一样,判了空跟没判一样
				ctype = m["Type"].(string)
			}
			switch ctype {
			case "Normal":
				if m["Path"] != nil {
					_, err := os.Stat(m["Path"].(string))
					if err == nil {
						cpath = m["Path"].(string)
					}
				}
				rules = make(map[string]Rule)
				for k, v := range m {
					if k[len(k)-1] == '2' && m[k[:len(k)-1]+"1"] != nil {
						gogo := strings.Split(v.(string), "=") // go 是关键字
						come := strings.Split(m[k[:len(k)-1]+"1"].(string), "=")
						if len(gogo) != 2 || len(come) != 2 {
							errF(RuleErrStr, nil)
							return
						}
						rules[k[:len(k)-1]] = Rule{ComeKey: come[0], ComeValue: come[1], GoKey: gogo[0], GoValue: gogo[1]}
					}
				}
			case "Upload":
				if !havePath(m) {
					return
				}
				if m["Max"] == nil {
					errF(RuleMissingErrStr, nil)
					return
				} else {
					cmax = m["Max"].(string)
				}
				if m["Upkey"] == nil {
					errF(RuleMissingErrStr, nil)
					return
				} else {
					cupkey = m["Upkey"].(string)
				}
			default:
				if !havePath(m) {
					return
				}
			}
		}

		// 判断这条请求发全了木有
		if cport == "" {
			if cip != "" {
				confIp = cip
				errF(SuccessErrStr, nil)
				return
			}
			if m["Start"] != nil && m["Start"].(string) == "ok" {
				fmt.Print("PortRules:")
				fmt.Println(PortRules)
				l := NewEasyListenToPortRuleList(confIp, PortRules...)
				l.Listen()
				go l.Commit() // Commit方法里面有个死锁,要是在主线程运行这个方法后面的都不用运行了
				errF(SuccessErrStr, nil)
				return
			}
		} else {
			mr := MuxRule{MuxPath: cmux}
			if ctype == "File" || ctype == "FileAll" || ctype == "Dir" {
				if ctype == "File" {
					mr.Type = IsAFile
				} else if ctype == "FileAll" {
					mr.Type = IsAFileAll
				} else if ctype == "Dir" {
					mr.Type = IsADir
				}
				mr.Rules = append(mr.Rules, Rule{ComeKey: "path", ComeValue: cpath})
				for _, p := range PortRules {
					if p.PortPath == cport {
						for _, m := range p.Muxs {
							if m.MuxPath == cmax {
								m.Type = mr.Type
								m.Rules = append(m.Rules, mr.Rules...)
								errF(SuccessErrStr, nil)
								return
							}
						}
						p.Muxs = append(p.Muxs, mr)
						errF(SuccessErrStr, nil)
						return
					}
				}
				p := PortRule{PortPath: cport, Muxs: []MuxRule{}}
				p.Muxs = append(p.Muxs, mr)
				PortRules = append(PortRules, p)
				errF(SuccessErrStr, nil)
				return
			} else if ctype == "Upload" {
				mr.Type = IsAReceive
				mr.Rules = append(mr.Rules, Rule{ComeKey: "path", ComeValue: cpath})
				mr.Rules = append(mr.Rules, Rule{ComeKey: "max", ComeValue: cmax, GoKey: "key", GoValue: cupkey})
				for _, p := range PortRules {
					if p.PortPath == cport {
						for _, m := range p.Muxs {
							if m.MuxPath == cmax {
								m.Type = mr.Type
								m.Rules = append(m.Rules, mr.Rules...)
								errF(SuccessErrStr, nil)
								return
							}
						}
						p.Muxs = append(p.Muxs, mr)
						errF(SuccessErrStr, nil)
						return

					}
				}
				p := PortRule{PortPath: cport, Muxs: []MuxRule{}}
				p.Muxs = append(p.Muxs, mr)
				PortRules = append(PortRules, p)
				errF(SuccessErrStr, nil)
				return
			} else if ctype == "Normal" {
				mr.Type = IsAFunc
				if cpath != "" {
					prs, err := readOneConfigFile(cpath)
					if err != nil {
						errF(SetErrStr, err)
					}
					PortRules = append(PortRules, prs...)
				}
				for _, r := range rules {
					mr.Rules = append(mr.Rules, r)
				}
				for _, p := range PortRules {
					if p.PortPath == cport {
						for _, m := range p.Muxs {
							if m.MuxPath == cmax {
								m.Type = mr.Type
								m.Rules = append(m.Rules, mr.Rules...)
								errF(SuccessErrStr, nil)
								return
							}
						}
						p.Muxs = append(p.Muxs, mr)
						errF(SuccessErrStr, nil)
						return

					}
				}
				p := PortRule{PortPath: cport, Muxs: []MuxRule{}}
				p.Muxs = append(p.Muxs, mr)
				PortRules = append(PortRules, p)
				errF(SuccessErrStr, nil)
				return
			}

		}

	}
}
