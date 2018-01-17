package ai

import(
	"encoding/json"
)

type PicWord struct{
	LogId int64 `json:"log_id"`
	WordResultNum int `json:"words_result_num"`
	WordsResults []Word `json:"words_result"`
}
type Word struct{
	Words string `json:"words"`
}

/*
  将AI文字识别的json结果解析成为题目和答案
  传入参数：json字符串
  返回参数：题目；答案切片；错误
*/
func GetQA(jsonStr string) (question string, answers []string, err error){
	buf := []byte(jsonStr)
	var picWord PicWord
	err = json.Unmarshal(buf, &picWord)
	if err != nil{
		return
	}
	qs := picWord.WordsResults[0:picWord.WordResultNum - 3]
	for _, v := range qs{
		question += v.Words
	}
	as := picWord.WordsResults[picWord.WordResultNum - 3:]
	for _, v := range as{
		answers = append(answers, v.Words)
	}
	return
}