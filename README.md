# EasyWeb
### 用户只需要简单配置一下端口号和路由地址即可完成建站，不依赖任何环境。服务器程序支持根据规则返回请求应得的返回值,文件上传，下载，文件或文件夹网络映射。如果需要进行逻辑处理也只需要关注逻辑部分即可，联网部分完全不用关心

### 更新信息,因为我不会做界面,所有把接口部分做好了,只要改动一下main文件里的启动方式,用GUI_WEB方式启动即可.你可以很轻易的将自己的界面和程序连接在一起,通过网络将他俩进行关联.控制接口可以设置规则和配置,也可以启动或关闭程序,但要注意关闭了程序你可能得手动在打开这个程序.

### 暂时还没有做完，但主要功能已经完成了,只要简单写一下配置文件即可,运行的时候使用-h或-help会获得参数提示。因为是新手，代码肯定写的不是很好，希望多多见谅，也希望这个项目能帮到需要搭建服务器但想偷懒的人们。

### 因为做开发或者建站的时候经常要写联网部分，但我很懒，所以做了这么个工具，工具分如下几个功能：

  <br><br>  1）对本机上的IP地址进行监听：本地回环地址，局域网地址，公网的IP地址都能监听（前提是你得有公网IP），可以根据需求监听某一个地址，也可以全部监听。
  <br><br>  2）可以同时监听一个或多个端口，只要简单的设置一下端口号即可，也可以根据情况分批监听端口号。暂时不支持停止监听，需要手动重启程序
  <br><br>  3）可以设置路由地址，可以每个端口号都可以设置任意多个要监听的路由地址。一定要先设置路由地址再启动端口监听。重启做出来应该就能解决这个问题了。
  <br><br>  4）目前支持功能：
  <br>            * 傻瓜式设置文件夹映射（将文件夹直接映射到网络上，用浏览器输入设置好的网址就能查看文件夹内容并可以下载其中的文件）
  <br>            * 傻瓜式设置上传路径（只要设置好你想设置的端口号和地址，并设置文件要储存的位置，选择上传模式即可。轻松向服务器传输文件，文件大小限制也可以设置）
  <br>            * 傻瓜式设置文件映射（可以将文件和你想要的地址进行绑定，网址中可以不出现文件名和后缀。支持将一个文件夹里的文件一次性映射出去（这种方法只能选择文件夹对应的路径，用户通过路径+文件名进行下载，文件夹本身不会被映射出去））
  <br><br>  5) 目前支持用户自定义逻辑，但得写代码处理，不过联网部分和数据解析部分处理完了，用户只要处理解析好的数据并返回要发送的结果即可。目前支持Post请求和Get请求，Post请求支持Json数据          
  <br><br>  6)当然,既然是用来偷懒的肯定不写代码也能用,是不过规则要在规则文件中写死.服务器会根据规则和用户发送的请求返回对应的数据
 
### 实际用途：

比如你想要搭建一个人网站，静态文件和网页文件直接映射出去，网站就搭建完了。就设置一下端口号和网址（IP地址可以不设置）

比如你想要搭建一个自用网盘，设置一个上传路径，在设置一个文件夹映射或者文件映射即可，界面同个人网站，简单干脆。

比如你做了个客户端想进行测试，但服务器端还没做好，只要设置一下规则，让服务器可以你想要的格式进行数据返回即可。完全不用重头写一个测试服务器软件有木有。

比如你和我一样懒，需要在一台新机器上跑一个server，但懒得装各种容器，也懒得配置环境，那恭喜你，你是最需要这个的人。

数据库还在做,但可以通过把数据记录到.easy文件中的方式实现储存和读取.具体怎么做很简单,自己想想就知道了.


# 懒人必备：EasyWeb！！！
  
  
