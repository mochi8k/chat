package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 予期せぬエラー
		panic(err.Error())
	} else {
		// 認証成功
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func errorChecker(err interface{}, provider common.Provider, message string) {
	if err != nil {
		log.Fatalln(message, provider, "-", err)
	}
}

// loginHandlerはサードパーティへのログインの処理を受け持つ.
// パスの形式: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {

	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]

	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		errorChecker(err, provider, "認証プロバイダーの取得に失敗:")

		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		errorChecker(err, provider, "GetBeginAuthURLの呼び出し中にエラーが発生:")

		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":

		// 認証プロバイダーの取得
		provider, err := gomniauth.Provider(provider)
		errorChecker(err, provider, "認証プロバイダーの取得に失敗:")

		// 認証
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		errorChecker(err, provider, "認証を完了できませんでした:")

		// ユーザーの取得
		user, err := provider.GetUser(creds)
		errorChecker(err, provider, "ユーザーの取得に失敗しました:")

		// クッキー生成
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
