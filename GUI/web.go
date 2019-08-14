package GUI

import "flag"

/*
我不会写界面,我也尝试过用拖得方式制作一个简易的界面,都失败了,
反正写了界面也是用发送请求的方式告诉程序配置信息,干脆你们自己写界面吧,自己做的界面用着也顺手些
 */
func WEB() {
	ip := flag.String("ip", "0.0.0.0", "指定要监听的端口,注意:这个端口是用来接收设置的")
	port := flag.String("port", "10086", "设置程序监听的端口,程序将从此端口获得参数信息")
	mux := flag.String("mux", "/", "设置程序监听的路由地址,程序将从端口的此路由获得参数信息")
	flag.Parse()
	confIp = *ip
	go JieXiAndJianTing(*port, *mux)
	select {} // 主线程空着备用挺好
}
