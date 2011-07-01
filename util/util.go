package util

import (
    "appengine"
    "os"
    "io"
    "json"
    "io/ioutil"
    "http"
    "fmt"
)

import (
    "oauth"
)

func ServeError(c appengine.Context, w http.ResponseWriter, err os.Error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "Internal Server Error")
	c.Errorf("%v", err)
}

func ReadToken(token interface{}, filename string) os.Error {
    b, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    return json.Unmarshal(b, token)
}

func WriteToken(token interface{}, filename string) os.Error {
    b, err := json.Marshal(token)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(filename, b, 0666)
}

func File(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    var token oauth.AccessToken
    /*
       token.Token = util.RequestToken
       token.Secret = util.RequestSecret
       err := writeToken(token, "sina.json")
       if err != nil {
           fmt.Fprintf(w, "writeToken Error %v\n", err)
           c.Infof("writeToken Error\n")
       }
       fmt.Fprintf(w, "writeToken OK<br/>")
    */
    err := ReadToken(&token, "sina.json")
    if err != nil {
        fmt.Fprintf(w, "readToken Error %v\n", err)
        c.Infof("readToken Error")
    }
    fmt.Fprintf(w, "readToken OK %v \n", token)
}
