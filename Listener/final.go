package Listener

import (
	"errors"
	"net/http"
)

// 路由地址及对应规则
type easyMux map[string]reFunc

// 路由类型
type MuxType int

// 文件类型
type FileType int

// 请求类型
type ResponseType int

// 规则
type reFunc struct {
	IsListenNow bool    // 是否正在被监听
	MuxType     MuxType // 规则类型
	RequestFunc func(http.ResponseWriter, *http.Request)
	Handle      http.Handler
}

// 一个网络监听器对象
type EasyListen struct {
	ip    string                    // 监听的IP地址 默认全部监听
	muxs  map[string]*http.ServeMux // 各端口号对应路由
	ports map[string]easyMux        // 监听的端口号和对应的路由规则；
}

const (
	IsAFile    MuxType = 1 // 文件映射型
	IsAFunc    MuxType = 2 // 处理数据型
	IsAReceive MuxType = 3 // 接收文件型
	IsADir     MuxType = 4 // 文件夹映射类型
)

const (
	AFIle FileType = 1 // 一个文件
	ADIR  FileType = 2 // 一个文件夹
)

const (
	EPOST ResponseType = 1 // post
	EGET  ResponseType = 2 // get
)

var (
	FileNotDirError error = errors.New("Want File But Dir")
	DirNotFileError error = errors.New("Want Dir But File")
)
