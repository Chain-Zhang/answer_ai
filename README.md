# answer_ai

[go实现西瓜视频花椒直播等平台智能答题](http://www.chairis.cn/blog/article/71)

下载代码：
```
git clone https://github.com/Chain-Zhang/answer_ai.git
```

# 介绍

最近出了很多答题平分奖金的直播，只要能够连续答对12道题，就能与所有答对的人一起平分奖池里的奖金，20万到500万不等。当这个时候，我才体会到“书到用时方恨少”这句至理名言。这时突然想到，咦！我们不是有无所不知的互联网吗，题目拿到百度中一搜不就完了。可是一看答题时间只有10秒，尽管我有着单身20多年的手速，也愣是做不到呀。再一想，我特么是程序猿呀，这种事还需要我亲自动手？

于是一通百度，找到了个大神的java智能答题的源码，这里把大神的源码地址贴出来供大家参考：https://github.com/lingfengsan/MillionHero

然而，我学了这么一大段时间的go语言，能不能用go来实现一下呢。
于是就动手尝试了一下，思路与前面提到的java的工具差不多。下面就来说道说道我是怎么实现的。

# 思路

1. 手机与电脑连接，并打开直播页面
2. 当页面出题时，通过adb截图并保存到电脑
3. 通过百度AI文字识别，提取图片中的题目和选项的文字
4. 使用百度搜索并，然后统计搜索得到结果数量
5. 比较搜索到的结果数量并排序
6. 否定的问题选择数量最少的选项，肯定的问题选数量最多的选项。

# 环境

## 硬件
* windows电脑一台
* 安卓手机一部
* 安卓数据线一根

## 软件
* golang 开发环境
* adb 安卓调试驱动

## 其他
* 百度AI开发者平台创建一个文字识别的应用

## 环境搭建
硬件就没有什么好说的了。这里主要说下软件。
### golang开发环境
首先，肯定是要下载安装包啦，这里给个下载地址，自己根据情况选择版本下载：[golang安装包](https://www.golangtc.com/download) (i386表示x86，amd64表示x64)。

安装完Go之后，我们最好还是检查一些所有的环境变量是否正常。主要的环境变量有以下几个：

* GOROOT：Go的安装目录
* GOPATH：用于存放Go语言Package的目录，这个目录不能在Go的安装目录中
* GOBIN：Go二进制文件存放目录，写成%GOROOT%\bin就好
* PATH：需要将%GOBIN%加在PATH变量的最后，方便在命令行下运行Go

完成之后在cmd窗口输入：`go version`
![go version](http://opgmvuzyu.bkt.clouddn.com/1516068224285.png)

如图所示，表示我们已经安装配置成功。

然后就是IDE了，这个就更简单了。直接用记事本都可以，当然也可以用些轻量的编辑器，vscode, vim都是可以的。也可以用goland等。这些看自己的爱好。反正我是用的vscode。

### adb安装
adb的全称为Android Debug Bridge 调试桥，是连接Android手机与PC端的桥梁，通过adb可以管理、操作模拟器和设备，如安装软件、查看设备软硬件参数、系统升级、运行shell命令等。

这里先给一个下载地址：[adb下载地址](http://download.csdn.net/download/zzceng/10204041) （有积分的大佬们从我这里下吧，我一分都没有了，想赚点分）

下载完成后安装好即可。然后把安装好的路径配置到环境变量中去，方便我们在cmd窗口下使用adb命令。配置好后，可以在cmd窗口下执行`adb devices` 命令：

![adb devices](http://opgmvuzyu.bkt.clouddn.com/1516068810427.png)

从图中可以看到，这里我们启动了adb，并且给了个设备列表，因为我没有连接安卓设备，所以没有东西显示。
这个时候，我们把安卓手机用数据线连接到电脑，并在手机上打开USB调试选项。`设置->开发者选项->USB调试`，不同的品牌的手机可能有差别，百度一下你就知道。

有时候可能做到这些还是列不出你的设备。这时候再需要做以下事情：
1. 在计算机管理中设备管理中找到你的设备，然后右击->属性->详细信息->在详细信息页面的属性中找到硬件ID，再复制的硬件ID，我的手机是魅族，我的硬件ID是：2A45
2. 在`C:\Users\你的用户名\.android`目录下找到adb_usb.ini文件，如果没有自己新建。然后把你刚刚复制的硬件ID写进去，由于这个ID是16进制的，所以前面加上0x，
即：`0x2A45`。
3. 重启adb，停止Adb：`adb kill-server`，启动adb：` adb start-server`。
完成这些应该就可以了。如果还是不行，请自行百度。

至此，我们的环境算是完成了。

# 实验
实验之前，肯定是下载源码喽，当然还有少不了的依赖包。
这里我用了个baidu-ai-sdk的包。
可以通过以下命令完成安装：
```
go get github.com/chenqinghe/baidu-ai-go-sdk
```
然后通过git下载我的源码：
```
git clone https://github.com/Chain-Zhang/answer_ai.git
```
我们先看下main函数的内容
```
func main(){
	for {
		var cmd string
		fmt.Printf("> ")
		fmt.Scan(&cmd)
		switch cmd{
		case "1":
			ai.Start()
		case "2":
            ai.ExeCommand("cmd", []string{"/c", "adb", "devices"})
		case "exit":
			os.Exit(1)
		}
	}
}
```
从代码中可以看到，在程序运行的时候会等待用户的输入。
1. 当输入 ` 1 ` 时会进行截图答题的操作。
2. 当输入 ` 2 ` 时会列出与电脑连接的设备
3. 当输入 ` exit ` 时会退出程序。

下面我们在cmd窗口中进入我们代码的目录，执行以下命令来运行我们的程序：
```
go run main.go
```
然后输入2，看下是否有设备连接：
![查看设备](http://opgmvuzyu.bkt.clouddn.com/1516029263072.png)

然后手机打开直播，当主播出题时，输入1回车，这里实验所以手机直接打开一张图片，手机界面如下图：
![手机界面](http://opgmvuzyu.bkt.clouddn.com/1516029493265.png)

经过一系列的分析后，返回以下结果：
![答题结果](http://opgmvuzyu.bkt.clouddn.com/1516029450262.png)

根据 `否定的问题选择数量最少的选项，肯定的问题选数量最多的选项` 所以这一题选择： 2-c哩c哩舞。
