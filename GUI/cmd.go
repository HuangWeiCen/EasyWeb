package GUI

import (
	"flag"
	"log"
	"fmt"
)

func CMD() {
	str := flag.String("conf", "", "配置文件的路径,文件正文以config开头的那个")
	ip := flag.String("ip", "0.0.0.0", "配置文件里设置的ip地址暂时是摆设,如果要指定IP的话还要在设一遍")
	flag.Parse()
	if *str == "" {
		log.Fatal("再懒也不能懒到配置文件也不传入啊")
	}
	rules, err := ReadConfigFile(*str)
	if err != nil {
		fmt.Println(err)
	}
	listen := NewEasyListenToPortRuleList(*ip, rules...)
	listen.Listen()
	listen.Commit()
}

