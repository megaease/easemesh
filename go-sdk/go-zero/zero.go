package zero

import (
	"net/http"

	"github.com/megaease/easemesh/go-sdk/stdlib"
)

// ServeDefault is the same with stdlib ServeDefault.
// The caller must call it to activate default agent.
var ServeDefault = stdlib.ServeDefault

// EaseMeshHandler wraps handler of go-zero as middleware.
func EaseMeshHandler(next http.HandlerFunc) http.HandlerFunc {
	return stdlib.DefaultAgent.WrapHandleFunc(next)
}
