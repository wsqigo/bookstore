package testcase

import (
	"fmt"
	"golang.org/x/net/context"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"
)

/*
文档链接：https://juejin.cn/post/7227828958988976185
*/

func makeRequest(url string) (string, error) {
	// 创建 http.Client 客户端实例
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}

	// 读取响应内容
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(body), nil
}

func TestNewTimer(t *testing.T) {
	// 1. 创建一个timer对象，等待5秒钟
	timeout := time.NewTimer(5 * time.Second)

	ch := make(chan string, 1)
	go func() {
		// 2. 这里我们简单模拟一个需要执行10秒的操作
		time.Sleep(10 * time.Second)
		ch <- "hello world"
	}()

	// 3. 在select语句钟处理超时时间 或者请求正常返回
	select {
	case tm := <-timeout.C:
		// 执行任务超时处理
		fmt.Println("操作超时", tm)
		return
	case result := <-ch:
		// 执行正常业务流程
		fmt.Println(result)
	}

	// 停止timer
	if !timeout.Stop() {
		<-timeout.C
	}

	// 操作执行完成
	fmt.Println("操作执行完成")
}

func TestRequestWithTimer(t *testing.T) {
	url := "https://baidu.com"
	// 设置超时时间为5秒
	timeout := 5 * time.Second
	// 创建一个计时器，等待超时
	timer := time.NewTimer(timeout)

	// 创建一个 channel，用于接受请求结果
	ch := make(chan string, 1)

	// 启动协程执行请求
	go func() {
		result, err := makeRequest(url)
		if err != nil {
			ch <- fmt.Sprintf("Error: %s", err.Error())
			return
		}
		ch <- result
	}()

	// 等待超时或者请求结果返回
	select {
	case result := <-ch:
		fmt.Println(result)
	case <-timer.C:
		fmt.Println("Request timed out")
	}
	// 请求完成后，停止定时器
	if !timer.Stop() {
		<-timer.C
	}
}

func TestTimerCtx(t *testing.T) {
	// 创建一个timerCtx，设置超时时间为3秒
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// 调用cancel函数，释放占用的资源
	defer cancel()

	// 开启一个协程执行任务
	ch := make(chan string, 1)
	go func() {
		// 模拟任务执行，休眠5秒
		time.Sleep(5 * time.Second)
		ch <- "hello world"
	}()

	// 在主协程中等待timerCtx超时或任务完成
	select {
	case <-ctx.Done():
		fmt.Println("timeout")
	case result := <-ch:
		fmt.Println(result)
	}
}

func TestRequestWithTimerCtx(t *testing.T) {
	url := "https://baidu.com"
	// 创建一个不带超时的context
	ctx := context.Background()

	// 1. 创建一个带超时的timerCtx
	timeout := 5 * time.Second
	timerCtx, cancel := context.WithTimeout(ctx, timeout)
	//5. 在函数返回时，调用取消函数 cancel()，释放占用的资源。
	defer cancel()

	// 创建一个 channel，用于接收请求的结果
	ch := make(chan string, 1)

	// 2. 将子上下文传递给需要进行超时控制的函数, 启动协程执行请求
	go func() {
		result, err := makeRequest(url)
		if err != nil {
			ch <- fmt.Sprintf("Error: %s", err.Error())
			return
		}
		ch <- result
	}()

	// 函数可以通过调用 context.Context 对象的 Done() 方法来判断是否超时。
	// 如果 Done() 方法返回的 channel 被关闭，则意味着已经超时，需要及时停止当前任务并返回。
	select {
	case result := <-ch:
		fmt.Println(result)
	case <-timerCtx.Done():
		fmt.Println("Request timed out")
	}
}

func TestGetTime(t *testing.T) {
	fmt.Println(time.Now().UnixMilli())
	fmt.Println(time.Now().AddDate(0, -3, 0).UnixMilli())
}

/*
timer原理文档地址：https://juejin.cn/post/7228417023725207608
*/

// 为什么需要超时控制
// 文档地址：https://juejin.cn/post/7230308610080702520
func getResource() (string, error) {
	conn, err := net.Dial("tcp", "https://www.baidu.com")
	if err != nil {
		return "", err
	}

	defer conn.Close()

	// 发送请求并等待响应
	_, err = conn.Write([]byte("GET /resource HTTP/1.1\r\nHost: example.com\r\n\r\n"))
	if err != nil {
		return "", err
	}
	resp, err := ioutil.ReadAll(conn)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
