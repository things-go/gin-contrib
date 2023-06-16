package signature

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	httpsign "github.com/things-go/http-signature-go"
)

func Test_Validator_DigestUsingShared(t *testing.T) {
	t.Run("validate digest", func(t *testing.T) {
		tests := []struct {
			name          string
			validator     httpsign.Validator
			body          []byte
			Nonce         string
			signingMethod httpsign.SigningMethod
			signingKey    []byte
			signBody      string
			wantErr       error
		}{
			{
				name:          "empty body",
				validator:     NewDigestUsingSharedValidator(),
				body:          []byte{},
				Nonce:         "http signature!!",
				signingMethod: httpsign.SigningMethodHmacSha256,
				signingKey:    []byte("http signing key!!"),
				signBody:      "hmac-sha256=u3rJnV7cWSJxCavgWv3+Kne6EyHE43MgK/Vozywrqd8=",
				wantErr:       nil,
			},
			{
				name:          "not empty body",
				validator:     NewDigestUsingSharedValidator(),
				body:          []byte("http signature!!"),
				Nonce:         "",
				signingMethod: httpsign.SigningMethodHmacSha256,
				signingKey:    []byte("http signing key!!"),
				signBody:      "hmac-sha256=u3rJnV7cWSJxCavgWv3+Kne6EyHE43MgK/Vozywrqd8=",
				wantErr:       nil,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/a", bytes.NewReader([]byte(tt.body)))

				r.Header.Set(httpsign.Digest, tt.signBody)
				r.Header.Set(httpsign.Nonce, tt.Nonce)
				err := tt.validator.Validate(r, &httpsign.Parameter{
					Method: tt.signingMethod,
					Key:    tt.signingKey,
				})
				require.NoError(t, err)

				require.Equal(t, tt.wantErr, err)
			})
		}
	})
}
