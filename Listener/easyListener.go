package Listener

import (
	"net/http"
	"io/ioutil"
	"log"
	"os"
	"fmt"
	"io"
	"encoding/json"
)

// 进行规则设置
func (this *EasyListen) Listen() {
	for p, mux := range this.ports {
		for m, ref := range mux {
			if ref.IsListenNow {
				return
			}
			switch ref.MuxType {
			case IsADir:
				this.muxs[p].Handle(m, ref.Handle)
				break
			default:
				this.muxs[p].HandleFunc(m, ref.RequestFunc)
				break
			}
			ref.IsListenNow = true
			fmt.Println("||		添加路由规则" + p + m)
		}
	}
	fmt.Println("")
}

// 启动服务器(先设置规则再启动服务器，否则规则不生效)
func (this *EasyListen) Commit() {
	for p, m := range this.muxs {
		fmt.Println("开始监听端口" + p)
		go http.ListenAndServe(this.ip+p, m)
	}
	select {}
}

// 添加一个要监听的端口
func (this *EasyListen) AddPort(port string) {
	port = ":" + port
	this.ports[port] = make(easyMux)
	this.muxs[port] = http.NewServeMux()
}

// 添加一个文件映射路径
func (this *EasyListen) AddFileMux(port, mux, filePath string) error {
	if err := wantFileType(filePath, AFIle); err != nil {
		return err
	}
	f := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Println("Read File Err:", err.Error())
		} else {
			log.Println("Send File:", filePath)
			w.Write(fileData)
		}
	}

	this.ports[port][mux] = reFunc{RequestFunc: f, IsListenNow: false, MuxType: IsAFile}
	return nil
}

// 添加一个文件夹的映射（只映射文件不映射文件夹本身和子文件夹）
func (this *EasyListen) AddDirAllFileMux(port, mux, dirPath string) error {
	if err := wantFileType(dirPath, ADIR); err != nil {
		if err == FileNotDirError {
			this.AddFileMux(port, mux, dirPath)
			return nil
		} else {
			return err
		}
	}
	f := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		des := dirPath + "/" + r.URL.Path[len(mux):len(r.URL.Path)]
		// string(os.PathSeparator)是反斜杠,文件地址都用的斜杠,输出出来怪怪的,干脆用"/"衔接算了
		fmt.Println("des:")
		fmt.Println("des:" + des)
		if err := wantFileType(des, AFIle); err == nil {
			if fileData, err := ioutil.ReadFile(des); err != nil {
				log.Println("Read File Err:", err.Error())
			} else {
				log.Println("Send File:", des)
				w.Write(fileData)
			}
		} else {
			if err == DirNotFileError {
				log.Println("File Is Dir", des)
			} else {
				log.Println("File Not Exit", des)
			}
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}
	mux = mux + "/"
	// go的路由匹配规则是末尾有"/"不管后面是啥都会匹配到这个路由.如果没有的话是完全匹配,就进不了这个路由了(可以删了试试效果)
	this.ports[port][mux] = reFunc{RequestFunc: f, IsListenNow: false, MuxType: IsAFileAll}
	return nil
}

// 添加一个文件夹的映射(连文件夹一起映射)
func (this *EasyListen) AddDirMux(port, mux, dirPath string) error {
	if err := wantFileType(dirPath, ADIR); err != nil {
		return err
	}
	handle := http.StripPrefix(mux, http.FileServer(http.Dir(dirPath)))
	this.ports[port][mux] = reFunc{Handle: handle, IsListenNow: false, MuxType: IsADir}
	return nil
}

// 添加一个接收文件用的路由(用post表单上传)
func (this *EasyListen) AddPostReceiveMux(maxMemory int64, port, postKey, mux, dirPath string) {
	f := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// 根据字段名获取表单文件
			err := r.ParseMultipartForm(maxMemory)
			if err != nil {
				log.Printf("Set max memory failed: %s\n", err)
			}
			formFile, header, err := r.FormFile(postKey)
			if err != nil {
				log.Printf("Get form file failed: %s\n", err)
				return
			}
			defer formFile.Close()

			// 创建保存文件
			destFile, err := os.Create(dirPath + "/" + header.Filename)
			if err != nil {
				log.Printf("Create failed: %s\n", err)
				return
			}
			defer destFile.Close()

			// 读取表单文件，写入保存文件
			_, err = io.Copy(destFile, formFile)
			if err != nil {
				log.Printf("Write file failed: %s\n", err)
				return
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	}
	this.ports[port][mux] = reFunc{RequestFunc: f, IsListenNow: false, MuxType: IsAReceive}
}

// 添加一个处理请求数据的路由
func (this *EasyListen) AddFuncMux(port, mux string, f func(w http.ResponseWriter, r *http.Request)) {
	this.ports[port][mux] = reFunc{RequestFunc: f, IsListenNow: false, MuxType: IsAFunc}
}

// 添加一个简单的请求处理路由(暂时只支持post 和 get)
func (this *EasyListen) AddEasyFuncMux(port, mux string, listener EasyHttpListener) {
	f := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		switch r.Method {
		case "POST":
			v := r.PostForm
			fmt.Println("vpost:")
			fmt.Println(v)
			w.Write(listener.EasyHttpListen(v, EPOST))
			break
		case "GET":
			v := r.URL.Query()
			fmt.Println("vget:")
			fmt.Println(v)
			w.Write(listener.EasyHttpListen(v, EGET))
		}

	}
	this.ports[port][mux] = reFunc{RequestFunc: f, IsListenNow: false, MuxType: IsAFunc}
}

// 添加一个简单的Json请求处理路由
func (this *EasyListen) AddEasyJsonFuncMux(port, mux string, listener EasyJSONListener) {
	f := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Method == "POST" {
			read, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Errorf("出现异常：%s", err)
				return
			}
			m := make(map[string]interface{})
			err = json.Unmarshal(read, &m)
			rm := listener.EasyJSONListen(m)
			rebt, err := json.Marshal(rm)
			if err != nil {
				fmt.Errorf("出现异常：%s", err)
				return
			}
			w.Write(rebt)
		}

	}
	this.ports[port][mux] = reFunc{RequestFunc: f, IsListenNow: false, MuxType: IsAFunc}
}
