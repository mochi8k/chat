package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"io/ioutil"
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
	"github.com/stretchr/objx"
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

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	// htmlからリクエスト内容を参照できる.
	t.templ.Execute(w, data)
}

func readProviders() func(name string) (string, string, string) {
	file, err := ioutil.ReadFile("./fixtures/provider.json")

	if err != nil {
		log.Fatal("JSON error")
	}

	var providers []provider
	json.Unmarshal(file, &providers)

	return func(name string) (string, string, string) {
		var id, secret, redirect string
		for _, provider := range providers {
			if provider.getName() == name {
				id = provider.getClientID()
				secret = provider.getClientSecret()
				redirect = "http://localhost:8080/auth/callback/" + name
			}
		}
		return id, secret, redirect
	}

}

func main() {
	var addr = flag.String("addr", ":8080", "アプリのドレス")

	// フラグを解釈
	flag.Parse()

	// set up gomniauth
	gomniauth.SetSecurityKey("tan security key")
	findProvider := readProviders()

	gomniauth.WithProviders(
		facebook.New(findProvider("facebook")),
		github.New(findProvider("github")),
		google.New(findProvider("google")),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1, // -1を指定することでクッキーを削除
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/room", r)
	go r.run()

	log.Println("Webサーバーを開始します。ポート: ", *addr)

	// start Web Server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
