package main

import (
    "fmt"
    "github.com/antchfx/xmlquery"
    "log"
    "net/http"
    "strings"
)

const xmlString = `
    <?xml version="1.0" encoding="utf-8"?>
        <users>
             <user ID="1">
                <name>TestUser</name>
                <password>test</password>
        </user>
    </users>
`

func checkAuth(name, password string) bool{
    document, err := xmlquery.Parse(strings.NewReader(xmlString))
    if err != nil {
        log.Fatal(err)
    }

    user := xmlquery.FindOne(document, fmt.Sprintf("/users/user[name/text()=%s and password/text()=%s]", name, password))

    return user != nil
}


func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("|  %s  |  %s  |  %s", r.Host, r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

func auth(next http.Handler) http.Handler {
    return http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            name := r.URL.Query().Get("name")
            password := r.URL.Query().Get("password")

            if name == "" && password != ""{
                w.WriteHeader(http.StatusBadRequest)
                w.Write([]byte("name or password is empty"))
                return
            }

            if ok := checkAuth(name, password); !ok{
                w.WriteHeader(http.StatusForbidden)
                w.Write([]byte("unknown user"))
                return
            }
            
            next.ServeHTTP(w, r)
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("you are an authorized user"))
}

func main() {
    http.Handle("/",  logging(auth(http.HandlerFunc(handler))))

    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}