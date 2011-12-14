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
    "strings"
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
    indexTemplate = template.Must(template.ParseFile("index.html"))
)

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
    http.HandleFunc("/get", get)
    http.HandleFunc("/file", util.File)
    http.HandleFunc("/sina", util.TokenSina)
    http.HandleFunc("/callback-sina", util.CallbackSina)
    http.HandleFunc("/PublicTimeLineSina", util.PublicTimeLineSina)
}

func get(w http.ResponseWriter, r *http.Request) {

    var httpRequest *http.Request

    c := appengine.NewContext(r)
    client := urlfetch.Client(c)

    url := r.FormValue("url")
    url = strings.TrimSpace(url)

    if url == "" {
        http.Redirect(w, r, "/", http.StatusFound)
        return
    }

    if len(url) > 5 {
        prefix := strings.ToLower(url[0:5])
        if !strings.HasPrefix(prefix, "http") && !strings.HasPrefix(prefix, "http") {
            url = fmt.Sprintf("http://%s", url)
        }

    } else {
        url = fmt.Sprintf("http://%s", url)
    }

    switch r.Method {
        default: {
            fmt.Fprintf(w, "Cannot handle method %v", r.Method)
            http.Error(w, "501 I only handle GET and POST", http.StatusNotImplemented)
            return
        }
        case "GET": {
            httpRequest, _ = http.NewRequest("GET", url, nil)
        }
        case "POST": {
            httpRequest, _ = http.NewRequest("POST", url, nil)
        }
    }

    httpRequest.Header.Set("Content-Type", "text/html; charset=utf-8")
    httpRequest.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux i686) AppleWebKit/534.24 (KHTML, like Gecko) Chrome/11.0.696.57 Safari/534.24,gzip(gfe)")

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

    for _, c := range req.Cookies() {
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
    _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Greeting", nil), &g)
    if err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusFound)
}
