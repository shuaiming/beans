package beans

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/shuaiming/beans/sessions"
	"github.com/shuaiming/openid"
)

const (
	// SesKeyOpenID Session key of OpenID
	SesKeyOpenID string = "openid_user"
	// SesKeyRedirect URL variable key for redirection after verified
	SesKeyRedirect string = "openid_redirect"
)

// OpenID login
type OpenID struct {
	prefix   string
	realm    string
	endpoint string
	openid   *openid.OpenID
}

// NewOpenID new OpenID
func NewOpenID(prefix, realm, endpoint string) *OpenID {
	return &OpenID{
		openid:   openid.New(realm),
		prefix:   prefix,
		realm:    realm,
		endpoint: endpoint,
	}
}

// ServeHTTPimp implement pod.Handler
func (o *OpenID) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if !strings.HasPrefix(r.URL.Path, o.prefix) {
		next(rw, r)
		return
	}

	if r.Method != "GET" && r.Method != "HEAD" {
		next(rw, r)
		return
	}

	s, ok := sessions.GetSession(r)
	if !ok {
		log.Printf("OpenID can not enabled without Sessions")
		next(rw, r)
		return
	}

	// verifyURL is url for OpenID Server back redirecetion
	verifyURL := fmt.Sprintf("%s/verify", o.prefix)

	switch r.URL.Path {

	case fmt.Sprintf("%s/login", o.prefix):

		// redirectURL is the url return back to after login finished
		// We will store it to "session" for later usage
		if redirectURL := r.URL.Query().Get(SesKeyRedirect); redirectURL != "" {
			s.Store(SesKeyRedirect, redirectURL)
		}

		// URL is url redirect to OpenID Server
		if authURL, err := o.openid.CheckIDSetup(o.endpoint, verifyURL); err == nil {
			http.Redirect(rw, r, authURL, http.StatusFound)
		} else {
			log.Println("OpenID error", err.Error())
		}

	case verifyURL:

		user, err := o.openid.IDRes(r)
		if err != nil {
			log.Println("OpenID error", err.Error())
			break
		}

		s.Store(SesKeyOpenID, user)

		if redirectURL, ok := s.Load(SesKeyRedirect); ok {
			http.Redirect(rw, r, redirectURL.(string), http.StatusFound)
		} else {
			http.Redirect(rw, r, o.realm, http.StatusFound)
		}

	default:
		next(rw, r)
	}
}

// GetOpenIDUser return User
func GetOpenIDUser(s sessions.Session) (map[string]string, bool) {
	user, ok := s.Load(SesKeyOpenID)

	if !ok {
		return nil, false
	}

	return user.(map[string]string), true
}
