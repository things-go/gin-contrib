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
	ErrFallback        func(c *gin.Context, statusCode int, err error)
	Hook               func(c *gin.Context, p *httpsign.Parameter)
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
		a.ErrFallback = func(c *gin.Context, statusCode int, err error) {
			c.String(statusCode, err.Error())
			c.Abort()
		}
	}
	return func(c *gin.Context) {
		if !a.SkipAuthentication(c) {
			parameter, err := a.Parser.ParseFromRequest(c.Request)
			if err != nil {
				a.ErrFallback(c, http.StatusBadRequest, err)
				return
			}
			err = a.Parser.Verify(c.Request, parameter)
			if err != nil {
				statusCode := http.StatusBadRequest
				if parameter.Scheme != httpsign.SchemeUnspecified &&
					errors.Is(err, httpsign.ErrSignatureInvalid) {
					if parameter.Scheme == httpsign.SchemeSignature {
						statusCode = http.StatusForbidden
					} else {
						statusCode = http.StatusUnauthorized
					}
				}
				a.ErrFallback(c, statusCode, err)
				return
			}
			if a.Hook != nil {
				a.Hook(c, parameter)
			}
		}
		c.Next()
	}
}
