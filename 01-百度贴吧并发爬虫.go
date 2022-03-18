package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func HttpGet2(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1 //将封装函数内部的错误，传出给调用者
		return
	}
	defer resp.Body.Close()
	//循环读取网页数据，传出给调用者
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		//累加每一次循环读到的buf数据，存入result；一次性返回。
		result += string(buf[:n])
	}
	return
}

//爬取单个页面的函数
func SpiderPage(i int, page chan int) {
	url := "https://tieba.baidu.com/p/7721471462?pn=" + strconv.Itoa(i)
	result, err := HttpGet2(url)
	if err != nil {
		fmt.Println("httpget err:", err)
		//continue //需要读取多个页，一个页面出错不用return，再读下一页就行了
		return
	}
	//fmt.Println("result=", result)
	//将读到的整网页数据，保存成一个文件
	f, err := os.Create("第" + strconv.Itoa(i) + "页" + ".html")
	if err != nil {
		fmt.Println("os.Create err:", err)
		//continue //需要读取多个页，一个页面出错不用return，再读下一页就行了
		return
	}
	f.WriteString(result)
	f.Close() //保存好一个文件，关闭一个文件
	page <- i //写数据
	//与主go程完成同步，协助其进行数据传递
}

//爬取页面操作
func working2(start, end int) {
	fmt.Printf("正在爬取第%d页到%d页...\n", start, end)
	//循环一次爬取一页
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第%d个页面爬取完成\n", <-page)
	}

}
func main() {
	//指定爬取起始，终止页
	var start, end int
	fmt.Println("请输入爬取的起始页（>=1）:")
	fmt.Scan(&start)
	fmt.Println("请输入爬取的终止页（>=strart）:")
	fmt.Scan(&end)
	working2(start, end)
}
