package golangtw

import (
    "appengine"
    "appengine/urlfetch"
    "appengine/datastore"
    "appengine/user"
    "http"
    "io/ioutil"
    "template"
    "time"
    "fmt"
)

import (
    "util"
)

type Greeting struct {
    Author  string
    Content string
    Date    datastore.Time
}

var (
    indexTemplate = template.MustParseFile("index.html", nil)
)

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
    http.HandleFunc("/get", get)
    http.HandleFunc("/file", util.File)
    http.HandleFunc("/sina", util.TokenSina)
    http.HandleFunc("/callback-sina", util.CallbackSina)
}

func get(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)

    url := r.FormValue("URL") + "/"

    httpRequest, _ := http.NewRequest("GET", url, nil)
    httpRequest.Header.Set("Content-Type", "text/html; charset=utf-8")
    httpRequest.UserAgent = "Mozilla/5.0 (X11; Linux i686) AppleWebKit/534.24 (KHTML, like Gecko) Chrome/11.0.696.57 Safari/534.24,gzip(gfe)"

    req, err := client.Do(httpRequest)
    if err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }

    c.Infof("HTTP Get %v returned status %v", url, req.Status)

    for k, v := range req.Header {
        for _, vv := range v {
            w.Header().Add(k, vv)
        }
    }

    for _, c := range req.SetCookie {
        w.Header().Add("Set-Cookie", c.Raw)
    }

    body, _ := ioutil.ReadAll(req.Body)
    defer req.Body.Close()
    fmt.Fprintf(w, "%v", string(body))
}

func root(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Greeting").Order("-Date").Limit(6)
    greetings := make([]Greeting, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }

    if err := indexTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
}

func sign(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    g := Greeting{
        Content: r.FormValue("content"),
        Date:    datastore.SecondsToTime(time.Seconds()),
    }
    if u := user.Current(c); u != nil {
        g.Author = u.String()
    }
    _, err := datastore.Put(c, datastore.NewIncompleteKey("Greeting"), &g)
    if err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusFound)
}
