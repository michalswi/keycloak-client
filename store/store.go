package store

import (
	"encoding/gob"

	"github.com/gorilla/sessions"
)

type token struct {
	Token string
}

var Store *sessions.CookieStore

func InitStore() error {
	Store = sessions.NewCookieStore([]byte("key-pairs"))
	gob.Register(token{})
	return nil
}

// var Store *sessions.FilesystemStore
// func InitStore() error {
// 	Store = sessions.NewFilesystemStore("", []byte("key-pairs"))
// 	gob.Register(map[string]interface{}{})
// 	return nil
// }
