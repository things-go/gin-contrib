package signature

import (
	"bytes"
	"io"
	"net/http"

	httpsign "github.com/things-go/http-signature-go"
)

type DigestUsingSharedValidator struct{}

// NewDigestValidator return pointer of new DigestValidator
func NewDigestUsingSharedValidator() *DigestUsingSharedValidator {
	return &DigestUsingSharedValidator{}
}

// Validate return error when checking digest match body
func (v *DigestUsingSharedValidator) Validate(r *http.Request, p *httpsign.Parameter) error {
	var err error
	var body []byte

	if r.ContentLength > 0 {
		// FIXME: using buffer to prevent using too much memory
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	} else {
		body = []byte(r.Header.Get(httpsign.Nonce))
	}
	headerDigest := r.Header.Get(httpsign.Digest)
	return httpsign.NewDigestUsingShared(p.Method).
		Verify(body, headerDigest, p.Key)
}
