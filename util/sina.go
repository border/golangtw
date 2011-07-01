package util

import (
    "appengine"
    "fmt"
    "io/ioutil"
    "http"
)

import (
    "oauth"
)

var provider = oauth.ServiceProvider{
    RequestTokenUrl:   "http://api.t.sina.com.cn/oauth/request_token",
    AuthorizeTokenUrl: "http://api.t.sina.com.cn/oauth/authorize",
    AccessTokenUrl:    "http://api.t.sina.com.cn/oauth/access_token",
}

var (
    callback string = "http://go.wifihack.net/callback-sina"
    //callback string = "http://127.0.0.1:8080/callback-sina"
)

var (
    RequestToken  string = ""
    RequestSecret string = ""
    AuthToken     string = ""
    AuthSecret    string = ""
    Code          string = ""
)


func TokenSina(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    var atoken oauth.AccessToken
    err := ReadToken(&atoken, "sina.json")
    if err != nil {
        ServeError(c, w, err)
        return
    }

    consumer := oauth.NewConsumer(c, atoken.Token, atoken.Secret, provider)
    atoken.Token = AuthToken
    atoken.Secret = AuthSecret
    if atoken.Token == "" || atoken.Secret == "" {
        c.Infof("Couldn't have auth token")
        var rtoken oauth.RequestToken
        rtoken.Token = RequestToken
        rtoken.Secret = RequestSecret
        if rtoken.Token == "" || rtoken.Secret == "" {
            c.Infof("Getting Request Token")
            rtoken, url, err := consumer.GetRequestTokenAndUrl("http://go.wifihack.net/")
            RequestToken = rtoken.Token
            RequestSecret = rtoken.Secret
            if err != nil {
                ServeError(c, w, err)
                return
            }
            c.Infof("Got rtoken: %v\n", rtoken)
            c.Infof("Visit this URL:", url)
            url = fmt.Sprintf("%s&oauth_callback=%s", url, callback)
            http.Redirect(w, r, url, http.StatusMovedPermanently)
            return
        }
    }
    http.Redirect(w, r, "/PublicTimeLineSina", http.StatusMovedPermanently)
}

func PublicTimeLineSina(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    var atoken oauth.AccessToken
    err := ReadToken(&atoken, "sina.json")
    if err != nil {
        ServeError(c, w, err)
        return
    }

    consumer := oauth.NewConsumer(c, atoken.Token, atoken.Secret, provider)

    //const url = "http://api.twitter.com/1/statuses/mentions.json"
    const url = "http://api.t.sina.com.cn/statuses/public_timeline.json"
    //const url = "http://api.t.sina.com.cn/account/verify_credentials.json"
    //const url = "http://api.t.sina.com.cn/statuses/user_timeline.json?user_id=1851069237"
    c.Infof("GET %v", url)
    atoken.Token = AuthToken
    atoken.Secret = AuthSecret
    resp, err := consumer.Get(url, nil, &atoken)
    if err != nil {
        ServeError(c, w, err)
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    fmt.Fprintf(w, "%v", string(body))
    return
}

func CallbackSina(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    Code = r.FormValue("oauth_verifier")
    c.Infof("[CallbackSina]RequestToken:%v, rTokenSecret: %v, Code: %v", RequestToken, RequestSecret, Code)

    var atoken oauth.AccessToken
    err := ReadToken(&atoken, "sina.json")
    if err != nil {
        ServeError(c, w, err)
        return
    }

    consumer := oauth.NewConsumer(c, atoken.Token, atoken.Secret, provider)
    var rtoken oauth.RequestToken

    rtoken.Token = RequestToken
    rtoken.Secret = RequestSecret
    if rtoken.Token != "" && rtoken.Secret != "" && Code != "" {
        tok, err := consumer.AuthorizeToken(&rtoken, Code)
        if err != nil {
            ServeError(c, w, err)
        }
        c.Infof("atoken: %v\n", tok)
        AuthToken = tok.Token
        AuthSecret = tok.Secret
        c.Infof("Get AuthToken: %v, AuthSecret: %v\n", AuthToken, AuthSecret)
        http.Redirect(w, r, "/PublicTimeLineSina", http.StatusMovedPermanently)
        return
    }
    http.Redirect(w, r, "/", http.StatusMovedPermanently)
    return
}
