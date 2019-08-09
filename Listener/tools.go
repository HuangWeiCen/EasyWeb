package Listener

import (
	"os"
	"net/http"
)

// 判断文件是文件夹还是文件，如果类型不匹配会返回错误
func wantFileType(filePath string, wantType FileType) error {
	desStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if desStat.IsDir() {
		if wantType == ADIR {
			return FileNotDirError
		}
	} else {
		if wantType == AFIle {
			return DirNotFileError
		}
	}
	return nil
}

// 传入第一个参数为监听的IP地址，往后都是端口号
func NewEasyListener(ip string) (listener *EasyListen) {
	p := make(map[string]easyMux)
	m := make(map[string]*http.ServeMux)
	listener = &EasyListen{}
	listener.ports = p
	listener.muxs = m
	listener.ip = ip
	return
}
