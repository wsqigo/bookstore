package errhdl

import (
	"bookstore/demo/web"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.RegisterError(http.StatusNotFound, []byte(`
<html>
	<body>
		<h1>NOT FOUND</h1>
	</body>
</html>
`)).
		RegisterError(http.StatusBadRequest, []byte(`
<html>
	<body>
		<h1>请求不对</h1>
	</body>
</html>
`))
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))
	server.Start(":8081")
}
