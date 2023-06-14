package signature

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	httpsign "github.com/things-go/http-signature-go"
)

const (
	readID       = httpsign.KeyId("read")
	readKey      = "1234"
	invalidKeyID = httpsign.KeyId("invalid key")
)

var (
	minimumRequiredHeaders = []string{"(request-target)", "date", "digest"}
	submitHeader           = []string{"(request-target)", "date", "digest"}
)

// mock interface always return true
type dateAlwaysValid struct{}

func (v *dateAlwaysValid) Validate(r *http.Request, _ *httpsign.Parameter) error { return nil }

// mock interface always return true
type timestampAlwaysValid struct{}

func (v *timestampAlwaysValid) ValidateTimestamp(t int64) error { return nil }

func runTest(
	signingMethod httpsign.SigningMethod,
	headers []string,
	v []httpsign.Validator,
	req *http.Request,
) *gin.Context {
	gin.SetMode(gin.TestMode)
	parser := httpsign.NewParser(
		httpsign.WithMinimumRequiredHeaders(headers),
		httpsign.WithValidators(v...),
		httpsign.WithValidatorCreated(&timestampAlwaysValid{}),
	)
	parser.RegisterSigningMethod(
		signingMethod.Alg(),
		func() httpsign.SigningMethod {
			return signingMethod
		},
	)
	_ = parser.AddMetadata(readID, httpsign.Metadata{
		Scheme: httpsign.SchemeUnspecified,
		Alg:    signingMethod.Alg(),
		Key:    []byte(readKey),
	})

	auth := Authenticator{
		Parser: parser,
	}
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = req
	auth.Authenticated()(c)
	return c
}

func TestAuthenticated_NoParser(t *testing.T) {
	require.Panics(t, func() {
		a := &Authenticator{}
		a.Authenticated()
	})
}

func TestAuthenticatedHeader_NoSignature(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	require.NoError(t, err)
	c := runTest(httpsign.SigningMethodHmacSha256, minimumRequiredHeaders, nil, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
}

func TestAuthenticatedHeader_WrongKey(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	require.NoError(t, err)

	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	param := httpsign.Parameter{
		KeyId:     invalidKeyID,
		Signature: "",
		Algorithm: "",
		Created:   0,
		Expires:   0,
		Headers:   submitHeader,
		Scheme:    httpsign.SchemeAuthentication,
		Method:    httpsign.SigningMethodHmacSha512,
		Key:       []byte(invalidKeyID),
	}
	err = param.MergerHeader(req)
	require.NoError(t, err)

	c := runTest(httpsign.SigningMethodHmacSha256, minimumRequiredHeaders, nil, req)
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())
}

func TestAuthenticate_SchemeAuthentication_InvalidSignature(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	require.NoError(t, err)

	req.Header.Set("Date", time.Date(1990, time.October, 20, 0, 0, 0, 0, time.UTC).Format(http.TimeFormat))
	req.Header.Set("digest", "xxxx")
	param := httpsign.Parameter{
		KeyId:     readID,
		Signature: "invalid signature",
		Algorithm: "",
		Created:   0,
		Expires:   0,
		Headers:   minimumRequiredHeaders,
		Scheme:    httpsign.SchemeAuthentication,
		Method:    httpsign.SigningMethodHmacSha512,
		Key:       []byte(readKey),
	}
	mockMergerHeader(&param, req)

	c := runTest(
		httpsign.SigningMethodHmacSha512,
		minimumRequiredHeaders,
		[]httpsign.Validator{&dateAlwaysValid{}},
		req,
	)
	assert.Equal(t, http.StatusUnauthorized, c.Writer.Status())
}

func TestAuthenticate_SchemeSignature_InvalidSignature(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	require.NoError(t, err)

	req.Header.Set("Date", time.Date(1990, time.October, 20, 0, 0, 0, 0, time.UTC).Format(http.TimeFormat))
	req.Header.Set("digest", "xxxx")
	param := httpsign.Parameter{
		KeyId:     readID,
		Signature: "invalid signature",
		Algorithm: "",
		Created:   0,
		Expires:   0,
		Headers:   minimumRequiredHeaders,
		Scheme:    httpsign.SchemeSignature,
		Method:    httpsign.SigningMethodHmacSha512,
		Key:       []byte(readKey),
	}
	mockMergerHeader(&param, req)

	c := runTest(
		httpsign.SigningMethodHmacSha512,
		minimumRequiredHeaders,
		[]httpsign.Validator{&dateAlwaysValid{}},
		req,
	)
	assert.Equal(t, http.StatusForbidden, c.Writer.Status())
}

func TestAuthenticate_Success(t *testing.T) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/", nil)
	require.NoError(t, err)

	req.Header.Set("Date", time.Date(1990, time.October, 20, 0, 0, 0, 0, time.UTC).Format(http.TimeFormat))
	req.Header.Set("digest", "xxxx")
	param := httpsign.Parameter{
		KeyId:     readID,
		Signature: "",
		Algorithm: "",
		Created:   0,
		Expires:   0,
		Headers:   minimumRequiredHeaders,
		Scheme:    httpsign.SchemeAuthentication,
		Method:    httpsign.SigningMethodHmacSha512,
		Key:       []byte(readKey),
	}
	err = param.MergerHeader(req)
	require.NoError(t, err)
	c := runTest(
		httpsign.SigningMethodHmacSha512,
		minimumRequiredHeaders,
		[]httpsign.Validator{&dateAlwaysValid{}},
		req,
	)
	assert.Equal(t, http.StatusOK, c.Writer.Status())
}

func mockMergerHeader(p *httpsign.Parameter, r *http.Request) {
	p.Algorithm = p.Method.Alg()

	b := strings.Builder{}
	hd := httpsign.SignatureHeader
	if p.Scheme == httpsign.SchemeAuthentication {
		hd = httpsign.AuthorizationHeader
		b.WriteString("Signature ")
	}
	b.WriteString(fmt.Sprintf(`keyId="%s",`, p.KeyId))
	b.WriteString(fmt.Sprintf(`algorithm="%s",`, p.Algorithm))
	if p.Created > 0 {
		b.WriteString(fmt.Sprintf(`created=%d,`, p.Created))
	}
	if p.Expires > 0 {
		b.WriteString(fmt.Sprintf(`expires=%d,`, p.Expires))
	}
	b.WriteString(fmt.Sprintf(`headers="%s",`, strings.Join(p.Headers, " ")))
	b.WriteString(fmt.Sprintf(`signature="%s"`, p.Signature))
	r.Header.Set(hd, b.String())
}
