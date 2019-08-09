package main

import (
	"./GUI"
	"fmt"
)

func main() {
	testReadFile()

}
func testReadFile() {
	rules, err := GUI.ReadConfigFile("C:/Users/huang/Desktop/GolangProject/EasyWeb/Demo/TestConf.easy")
	if err != nil {
		fmt.Println(err)
	}
	listen := GUI.NewEasyListenToPortRuleList("0.0.0.0", rules...)
	listen.Listen()
	listen.Commit()

}
