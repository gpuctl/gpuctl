package webapi_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gpuctl/gpuctl/internal/authentication"
	"github.com/gpuctl/gpuctl/internal/webapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

// Don't worry, this was generated for testing purposes :)
const demoPrivKey = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACBjK4r4g7acz0NMUdHdZR7IMdcbBykKLCjVGZA4af41ygAAAJDlmrP05Zqz
9AAAAAtzc2gtZWQyNTUxOQAAACBjK4r4g7acz0NMUdHdZR7IMdcbBykKLCjVGZA4af41yg
AAAEDFqhoLXMIqj+C810RJ2oUHLczGxXE9kneJse9y/LeNiWMriviDtpzPQ0xR0d1lHsgx
1xsHKQosKNUZkDhp/jXKAAAADERlbW8gU1NIIEtleQE=
-----END OPENSSH PRIVATE KEY-----`

var emptyAuthCookie = &http.Cookie{Name: authentication.TokenCookieName, Value: ""}

// TODO: Figure out how to test the happy path somehow

func TestOnboardNoKey(t *testing.T) {
	t.Parallel()

	serv := webapi.NewServer(nil, alwaysAuth{}, webapi.OnboardConf{
		DataDir: "/foo",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/admin/add_workstation", strings.NewReader(`{"hostname": "foo.net"}`))
	req.AddCookie(emptyAuthCookie)
	w := httptest.NewRecorder()

	serv.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "no ssh key")
}

func TestOnboardNoHostname(t *testing.T) {
	t.Parallel()

	sign, err := ssh.ParsePrivateKey([]byte(demoPrivKey))
	require.NoError(t, err)

	serv := webapi.NewServer(nil, alwaysAuth{}, webapi.OnboardConf{
		DataDir:  "/foo",
		Username: "root",
		Signer:   sign,
	})

	req := httptest.NewRequest(http.MethodPost, "/api/admin/add_workstation", strings.NewReader("{}"))
	req.Header.Add("Authorization", "Bearer 123")
	req.AddCookie(emptyAuthCookie)
	w := httptest.NewRecorder()

	serv.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), "hostname cannot be empty")
}
