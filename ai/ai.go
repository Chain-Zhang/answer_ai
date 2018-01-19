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
		pairs_baidu := make(PairList, len(a))
		pairs_sougou := make(PairList, len(a))
		pairs_360 := make(PairList, len(a))
		chanbaidu := make(chan Pair, len(a))
		chansougou := make(chan Pair, len(a))
		chan360 := make(chan Pair, len(a))
		for k, v := range a{
			fmt.Println("搜索关键字：", q + " " + v)
			go SeachBaidu(k, q, v, chanbaidu)
			go SeachSougou(k, q, v, chansougou)
			go Seach360(k, q, v, chan360)
		}
		for i := 0; i < len(a); i++{
            pairs_baidu[i] = <-chanbaidu
		}
		for i := 0; i < len(a); i++{
            pairs_sougou[i] = <-chansougou
		}
		for i := 0; i < len(a); i++{
            pairs_360[i] = <-chan360
		}
		
		sort.Sort(pairs_sougou)
		sort.Sort(pairs_baidu)
		sort.Sort(pairs_360)
		fmt.Println("\n\n百度搜索结果：")
		for _, v := range pairs_baidu{
            fmt.Printf("%s: %d \n", v.key, v.value)
		}
		fmt.Println("\n\n搜狗搜索结果：")
		for _, v := range pairs_sougou{
            fmt.Printf("%s: %d \n", v.key, v.value)
		}
		fmt.Println("\n\n360搜索结果：")
		for _, v := range pairs_360{
            fmt.Printf("%s: %d \n", v.key, v.value)
		}
	}else{
		fmt.Println("获取图片失败")
	}
}