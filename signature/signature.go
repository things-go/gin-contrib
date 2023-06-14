package signature

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	httpsign "github.com/things-go/http-signature-go"
)

// Authenticator is the gin authenticator middleware.
type Authenticator struct {
	Parser             *httpsign.Parser
	SkipAuthentication func(c *gin.Context) bool
	ErrFallback        func(*gin.Context, httpsign.Scheme, error)
}

// Authenticated returns a gin middleware which permits given permissions in parameter.
func (a *Authenticator) Authenticated() gin.HandlerFunc {
	if a.Parser == nil {
		panic("http sign parse must be initialized.")
	}
	if a.SkipAuthentication == nil {
		a.SkipAuthentication = func(c *gin.Context) bool { return false }
	}
	if a.ErrFallback == nil {
		a.ErrFallback = func(c *gin.Context, scheme httpsign.Scheme, err error) {
			statsCode := http.StatusBadRequest
			if scheme != httpsign.SchemeUnspecified &&
				errors.Is(err, httpsign.ErrSignatureInvalid) {
				if scheme == httpsign.SchemeSignature {
					statsCode = http.StatusForbidden
				} else {
					statsCode = http.StatusUnauthorized
				}
			}
			c.String(statsCode, err.Error())
		}
	}
	return func(c *gin.Context) {
		if !a.SkipAuthentication(c) {
			scheme, err := a.Parser.Parse(c.Request)
			if err != nil {
				a.ErrFallback(c, scheme, err)
				return
			}
		}
		c.Next()
	}
}
