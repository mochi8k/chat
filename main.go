package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/mochi8k/chat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	// htmlからリクエスト内容を参照できる.
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリのドレス")

	// フラグを解釈
	flag.Parse()

	// set up gomniauth
	gomniauth.SetSecurityKey("tan security key")
	gomniauth.WithProviders(
		facebook.New(
			"1659131477747431",
			"05d8695a82578f5bdb4e5bc58893d007",
			"http://localhost:8080/auth/callback/facebook",
		),
		github.New(
			"25799018a8946414a78b",
			"978387a957b2c4571a5a002333bb6c7405d5c385",
			"http://localhost:8080/auth/callback/github",
		),
		google.New(
			"219682017207-pa5443eapqatgbd6j3klpg7eutfij626.apps.googleusercontent.com",
			"FwDuL2xoeTCJxpi6uyXmQKgK",
			"http://localhost:8080/auth/callback/google",
		),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()

	log.Println("Webサーバーを開始します。ポート: ", *addr)

	// start Web Server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
