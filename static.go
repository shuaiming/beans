package beans

import (
	"net/http"
	"strings"
)

// Static serve static files
type Static struct {
	Prefix string
	Dir    http.FileSystem
	Index  bool
}

// NewStatic new Static
func NewStatic(prefix string, fs http.FileSystem, index bool) *Static {

	return &Static{
		Prefix: prefix, // statics url prefix
		Dir:    fs,
		Index:  index,
	}
}

// ServeHTTP implement pod.Handler
func (s *Static) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if !strings.HasPrefix(r.URL.Path, s.Prefix) {
		next(rw, r)
		return
	}

	// do not list dir when needed
	if !s.Index && strings.HasSuffix(r.URL.Path, "/") {
		http.Error(rw, "403 Forbidden", http.StatusForbidden)
		return
	}

	// fix http.StripPrefix dirList() bug
	if r.URL.Path == s.Prefix && !strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(rw, r, r.URL.Path+"/", http.StatusMovedPermanently)
		return
	}

	http.StripPrefix(s.Prefix, http.FileServer(s.Dir)).ServeHTTP(rw, r)
}
