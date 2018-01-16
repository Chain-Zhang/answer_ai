package ai

import (
	"fmt"
	"sort"
)

type Pair struct{
	key string
	value int64
}
type PairList []Pair
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].value < p[j].value }

/*
  开始作答。
*/
func Start(){
	path, b := GetPicPath()
	if b{
		fmt.Println("pic path: ", path)
		str, err := GetStringByBaiduai(path)
		if err != nil{
			fmt.Println("图片提取文字错误：", err.Error())
			return
		}
	    q, a, err := GetQA(str)
		if err != nil{
			fmt.Println("获取问题和答案失败：", err.Error())
			return
		}
		pairs := make(PairList, len(a))
		for k, v := range a{
			key := q + " " + v
			count := Seach(key)
			pairs[k] = Pair{key:fmt.Sprintf("%d=>%s", k + 1, v), value:count}
		}
		sort.Sort(pairs)
		fmt.Println("搜索结果：")
		for _, v := range pairs{
            fmt.Printf("%s: %d \n", v.key, v.value)
		}
	}else{
		fmt.Println("获取图片失败")
	}
}