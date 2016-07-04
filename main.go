package main

import (
	"log"
  "net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
      <html>
        <head>
          <title>チャット</title>
        <head>
        <body>
          Hello, Chat!
        </body>
      </html>
    `))
	})

	// start Web Server
  if err := http.ListenAndServe(":8000", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
