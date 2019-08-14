package GUI

import "errors"

var (
	NotKnowFileError error = errors.New("Not know file")
	NotRuleError     error = errors.New("Is config not rule")
	IsRuleError      error = errors.New("Is rule not config")
	WrongFormatError error = errors.New("Wrong format")
)

var (
	confIp    string     = "0.0.0.0"
	PortRules []PortRule = []PortRule{}
)

// 一条规则
type Rule struct {
	ComeKey   string      // 入参键
	ComeValue interface{} // 入参值
	GoKey     string      // 出参键
	GoValue   interface{} // 出参值
}

// 一条路径
type MuxRule struct {
	MuxPath string
	Rules   []Rule
	Type    MuxType
}

// 一个端口
type PortRule struct {
	Muxs     []MuxRule
	PortPath string
}

// 创建的监听器对象
type postListener struct {
	MuxRule
}

type MuxType int

type SetErrorString string

const (
	SetErrStr         SetErrorString = "设置期间程序抛出异常"
	SetFailErrStr     SetErrorString = "设置失败了"
	RuleMissingErrStr SetErrorString = "规则没有传全"
	RuleErrStr        SetErrorString = "传来的规则有问题"
	SuccessErrStr     SetErrorString = "成功了没异常"
)

const (
	IsAFile     MuxType = 1 // 文件映射型
	IsAFunc     MuxType = 2 // 处理数据型
	IsAReceive  MuxType = 3 // 接收文件型
	IsADir      MuxType = 4 // 文件夹映射类型
	IsAFileAll  MuxType = 5 // 映射一堆文件型
	IsAFileFunc MuxType = 6 // 映射一堆文件型
	IsAEasy     MuxType = 7 // 特殊类型（除了2，6都是特殊类型）
)
