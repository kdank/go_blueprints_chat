package main

import (
	"blueprints/trace"
	"flag"
	"github.com/stretchr/signature"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// templ represents a single template
type templateHandler struct {
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.templ == nil {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	}

	data := map[string]interface{}{
		"Host": r.Host,
	}
	log.Println("Host:", r.Host)

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
		log.Println("UserData:", data)
	}

	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application")
	flag.Parse() // parse the flag

	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		google.New("605169096484-hvgf6b5u60cdv0krm7r106b1v4lsuvk9.apps.googleusercontent.com",
			"0F1iTCjapfcoFUupsINMc-6X", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout)

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/chat",
		MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login",
		&templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
