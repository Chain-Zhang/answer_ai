package ai

import (
	"fmt"
	
	"io"
	"bufio"
	"os"
	"os/exec"

	"strconv"
	"strings"
	"time"

	"image"
	"image/png"

	"github.com/chenqinghe/baidu-ai-go-sdk/version/ocr"
	"answer_ai/conf"
)

var(
	cmdName = "cmd"
)

/*
  使用adb命令让手机截屏，并保存在sdcard里
*/
func getPic() bool{
	params := []string{"/c", "adb", "shell", "/system/bin/screencap", "-p", "/sdcard/screenshot.png"}
	return ExeCommand(cmdName, params)
}

/*
  使用adb命令将截图所得图片保存到电脑
  返回参数：电脑端图片路径；是否保存成功
*/
func savePic()(string, bool){
	timestamp := time.Now().Unix()
	picPath := conf.PicDirectory + strconv.FormatInt(timestamp, 10) + ".png"
    params := []string{"/c", "adb", "pull", "/sdcard/screenshot.png", picPath}
	b := ExeCommand(cmdName, params)
	if b{
		if pathExists(picPath){
			return picPath, true
		}
	}
	return "", false
}

/*
   判断文件是否存在
*/
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

/*
  完成手机截屏，并保存到电脑端的操作
  返回参数：电脑端图片路径；是否成功
*/
func GetPicPath()(string, bool){
	if getPic(){
		return savePic()
	}
	return "",false
}

/*
  通过百度AI从图片中识别文字
  传入参数：图片路径
  返回参数：识别的内容；错误
*/
func GetStringByBaiduai(picPath string)(string, error){
	picPath,err := cut(picPath)
	client := ocr.NewOCRClient(conf.BaiduAiAppKey, conf.BaiduAiSecretKey)
    fmt.Println(picPath)
	f, err := os.OpenFile(picPath, os.O_RDONLY, 0777)
	if err != nil {
		return "",err
	}
	rs, err := client.GeneralRecognizeBasic(f)
	if err != nil {
		return "",err
	}
	return string(rs), nil
}

/*
  裁剪掉图片中除题目和答案以外的文字
  传入参数：原图片路径
  传出参数：裁剪后的图片路径；错误
*/
func cut(path string) (picPath string, err error){
	reader, err := os.OpenFile(path, os.O_RDONLY, 0777)
	if err != nil {
		return
	}
	picPath = conf.PicDirectory + "cut.png"
	m, _, _ := image.Decode(reader)
	rgbImg := m.(*image.NRGBA)
	subImg := rgbImg.SubImage(image.Rect(conf.X1,conf.Y1,conf.X2,conf.Y2)).(*image.NRGBA)
	f, err := os.Create(picPath)
	defer f.Close()
	err = png.Encode(f, subImg)
	return
}

/*
  执行Dos命令
  传入参数：命令名称，一般为cmd；命令参数
  传出参数：是否执行成功
*/
func ExeCommand(cmdName string, params []string) bool{
	cmd := exec.Command(cmdName, params...)
	fmt.Printf("执行命令：%s \n", strings.Join(cmd.Args[2:], " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("error => ", err.Error())
		return false
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF{
			break
		}
		fmt.Println(line)
	}
	cmd.Wait()
	return true
}
