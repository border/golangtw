package hello

import (
    "fmt"
    "http"
    "template"
)

func init() {
    http.HandleFunc("/", root)
    http.HandleFunc("/sign", sign)
}

func root(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, guestbookForm)
}

const guestbookForm = `
<html>
  <body>
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`

func sign(w http.ResponseWriter, r *http.Request) {
    err := signTemplate.Execute(w, r.FormValue("content"))
    if err != nil {
        http.Error(w, err.String(), http.StatusInternalServerError)
    }
}

var signTemplate = template.MustParse(signTemplateHTML, nil)

const signTemplateHTML = `
<html>
  <body>
    <p>You wrote:</p>
    <pre>{@|html}</pre>
  </body>
</html>
`
