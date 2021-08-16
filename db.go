package beans

import (
	"context"
	"database/sql"
	"net/http"
)

// CtxKeyDB 用来标识Context中存放的值。
const CtxKeyDB key = "beans.DB"

// DB 读写数据库
type DB struct {
	*sql.DB
}

// NewDB new DB
func NewDB(db *sql.DB) *DB {
	return &DB{db}
}

// ServeHTTPimp implement pod.Handler
func (db *DB) ServeHTTP(
	rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	ctx := context.WithValue(r.Context(), CtxKeyDB, db.DB)
	next(rw, r.WithContext(ctx))

}

// GetDB return *sql.DB
func GetDB(r *http.Request) (*sql.DB, bool) {
	db := r.Context().Value(CtxKeyDB)
	if db == nil {
		return nil, false
	}

	return db.(*sql.DB), true
}
