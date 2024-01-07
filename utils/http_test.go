package utils

import (
	"fmt"
	"github.com/dpwgc/easierweb"
	"net/http"
	"testing"
	"time"
)

func TestHTTPClient(t *testing.T) {

	fmt.Println("\n[TestHTTPClient] start")

	router := easierweb.New()
	router.Any("/http/client/test", hello)
	go func() {
		err := router.Run(":80")
		if err != nil && err.Error() != "http: Server closed" {
			panic(err)
		}
	}()

	time.Sleep(3 * time.Second)

	_, result, err := GET("http://127.0.0.1/http/client/test", easierweb.Params{
		"User-Name": "easierweb",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[HTTPClientTest](GET)", string(result))
	time.Sleep(1 * time.Second)

	_, result, err = HEAD("http://127.0.0.1/http/client/test", easierweb.Params{
		"User-Name": "easierweb",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[HTTPClientTest](HEAD)", string(result))
	time.Sleep(1 * time.Second)

	_, result, err = OPTIONS("http://127.0.0.1/http/client/test", easierweb.Params{
		"User-Name": "easierweb",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[HTTPClientTest](OPTIONS)", string(result))
	time.Sleep(1 * time.Second)

	_, result, err = POST("http://127.0.0.1/http/client/test", []byte("test"), easierweb.Params{
		"User-Name": "easierweb",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[HTTPClientTest](POST)", string(result))
	time.Sleep(1 * time.Second)

	_, result, err = PUT("http://127.0.0.1/http/client/test", []byte("test"), easierweb.Params{
		"User-Name": "easierweb",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[HTTPClientTest](PUT)", string(result))
	time.Sleep(1 * time.Second)

	_, result, err = PATCH("http://127.0.0.1/http/client/test", []byte("test"), easierweb.Params{
		"User-Name": "easierweb",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[HTTPClientTest](PATCH)", string(result))
	time.Sleep(1 * time.Second)

	_, result, err = DELETE("http://127.0.0.1/http/client/test", easierweb.Params{
		"User-Name": "easierweb",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[HTTPClientTest](DELETE)", string(result))
	time.Sleep(1 * time.Second)

	err = router.Close()
	if err != nil {
		panic(err)
	}

	fmt.Println("\n[TestHTTPClient] end")
}

func hello(ctx *easierweb.Context) {
	ctx.WriteString(http.StatusOK, "hello "+ctx.Header.Get("User-Name"))
}
