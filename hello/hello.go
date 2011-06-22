package hello

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

type Greeting struct {
    Author  string
    Content string
    Date    datastore.Time
}

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
    http.HandleFunc("/get", get)
    http.HandleFunc("/google4a30971798e895bc.html", google)
}

func google(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "google-site-verification: google4a30971798e895bc.html")
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
    q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
    greetings := make([]Greeting, 0, 10)
    if _, err := q.GetAll(c, &greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }
    if err := guestbookTemplate.Execute(w, greetings); err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
}

var guestbookTemplate = template.MustParse(guestbookTemplateHTML, nil)

const guestbookTemplateHTML = `
<html>
  <body>
    {.repeated section @}
      {.section Author}
        <p><b>{@|html}</b> wrote:</p>
      {.or}
        <p>An anonymous person wrote:</p>
      {.end}
      <pre>{Content|html}</pre>
    {.end}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
    <form action="/get" method="post">
      <div><input type="text" name="URL" size="50"></input></div>
      <div><input type="submit" value="Http Proxy"></div>
    </form>
    <p></p>
    <div><a href="http://code.google.com/appengine/docs/go/">
    <img src="http://code.google.com/appengine/images/appengine-silver-120x30.gif" 
    alt="Powered by Google App Engine" />
    </a>
    </div>

  </body>
</html>
`

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
