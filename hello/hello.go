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
}

func get(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    client := urlfetch.Client(c)

    httpRequest, _ := http.NewRequest("GET", "http://www.sina.com.cn", nil)
    httpRequest.Header.Set("Content-Type", "text/html; charset=utf-8")

    req, err := client.Do(httpRequest)
    if err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
        return
    }
    for k, v := range req.Header {
        for _, vv := range v {
            w.Header().Add(k, vv)
            fmt.Printf("Key:%v, Value:%v\n", k, vv)
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
