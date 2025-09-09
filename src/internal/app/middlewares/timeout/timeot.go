package timeout

import (
	"net/http"
	"time"
)

type Middleware struct {}

func (rm *Middleware) Handle(next http.Handler) http.Handler {
	return http.TimeoutHandler(next, 5*time.Second, `{"error":"timeout"}`+"\n")
}