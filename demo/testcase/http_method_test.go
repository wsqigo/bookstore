package testcase

import (
	"fmt"
	"io"
	"net/http"
)

func readBodyOnce(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "read body failed: %v", err)
		return
	}

	// 类型转换，将 []byte 转化为 string
	fmt.Fprintf(w, "read the data: %s \n", string(body))

	// 尝试再次读取，啥也读不到，但是也不会报错
	body, err = io.ReadAll(r.Body)
	if err != nil {
		// 不会进来这里
		fmt.Fprintf(w, "read the data one more time got error: %v", err)
		return
	}
	fmt.Fprintf(w, "read the data one more time: [%s] and read data length %d", string(body), len(body))
}

func getBodyIsNil(w http.ResponseWriter, r *http.Request) {
	if r.GetBody == nil {
		fmt.Fprintf(w, "GetBody is nil \n")
	} else {
		fmt.Fprintf(w, "GetBOdy is not nil \n")
	}
}

func queryParams(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	fmt.Fprintf(w, "query is %v\n", values)
}

func form(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "before parse form %v\n", r.Form)
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "parse form error %v\n", r.Form)
		return
	}
	fmt.Fprintf(w, "after parse form %v", r.Form)
}
