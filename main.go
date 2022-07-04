/*
Serve is a very simple static file server in go
Usage:
	-p="8100": port to serve on
	-d=".":    the directory of static files to host
Navigating to http://localhost:8100 will display the index.html or directory
listing file.
*/
package main

import (
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"flag"
	"log"
	"net/http"
)

//go:embed resources/login.html
var loginPage string

//go:embed resources/fake.html
var fakePage string

func main() {
	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host. Default to current dir")
	secret := flag.String("secret", "", "secret to browse file. If provided, would check this secret cookie")
	fake := flag.String("fake", "", "if provided, serve a fake webpage instead")

	flag.Parse()

	var hashedSecret string
	if len(*secret) > 0 {
		hashedSecret = encode(*secret)
	}

	filehandler := http.FileServer(http.Dir(*directory))
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if len(hashedSecret) > 0 {
			// check cookie
			c, err := request.Cookie("secret")
			if err == http.ErrNoCookie || c.Value != hashedSecret {
				// handle login
				http.Redirect(writer, request, "/login", http.StatusTemporaryRedirect)
				return
			}
		}

		filehandler.ServeHTTP(writer, request)
	})
	http.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			writer.Write([]byte(loginPage))
			return
		} else {
			// assume POST otherwise
			err := request.ParseForm()
			if err != nil {
				log.Println("parse form error: ", err)
				return
			}

			code := request.FormValue("secret")
			if code == *fake {
				log.Println("serving page")
				writer.Write([]byte(fakePage))
				return
			}
			encoded := encode(code)
			if encoded == hashedSecret {
				// correct, set cookie
				cookie := &http.Cookie{
					Name:     "secret",
					Value:    encoded,
					Secure:   false,
					HttpOnly: false,
				}
				http.SetCookie(writer, cookie)
				http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
			} else {
				log.Println("invalid code input: ", code)
				writer.Write([]byte(loginPage))
			}
		}
	})

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func encode(secret string) string {
	encoder := sha256.New()
	encoder.Write([]byte(secret))
	return hex.EncodeToString(encoder.Sum(nil))
}
