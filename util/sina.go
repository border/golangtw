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
    callback string = "http://127.0.0.1:8080/callback-sina"
    )
/*
var (
    RequestToken  string = "d551777285524f8abca1347eebf00650"
    RequestSecret string = "9772bc3d3d2daa27f99fc6df61cb9a50"
    AuthToken     string = "4979774f91e92ef7e3c8483dc8299684"
    AuthSecret    string = "1f5a4aacf1a8f37c03dcba94752d65ba"
    Code          string = "226496"
)
*/

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
            c.Infof("rtoken: %v\n", rtoken)
            c.Infof("url: %v\n", url)
            if err != nil {
                ServeError(c, w, err) 
                return
            }
            c.Infof("Visit this URL:", url)
            url = fmt.Sprintf("%s&oauth_callback=%s", url, callback)
            //fmt.Fprintf(w, "Visit this URL %v<br>", url)
            http.Redirect(w, r, url, http.StatusMovedPermanently)
            return
        }

        c.Infof("Getting Access Token")
        if Code == "" {
            c.Infof("You must supply a -code parameter to get an Access Token.")
            return
        }
        tok, err := consumer.AuthorizeToken(&rtoken, Code)
        if err != nil {
            fmt.Fprintf(w, "%v", err)
        }
        c.Infof("atoken: %v\n", tok)
        atoken = *tok
    }

    PublicTimeLineSina(w, r, &atoken)
}

func aTokenSina(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

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
        AuthToken = atoken.Token
        AuthSecret = atoken.Secret 
        PublicTimeLineSina(w, r, tok)
    }
}
func PublicTimeLineSina(w http.ResponseWriter, r *http.Request, atoken *oauth.AccessToken) {
    c := appengine.NewContext(r)

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
    resp, err := consumer.Get(url, nil, atoken)
    if err != nil {
        fmt.Fprintf(w, "%v", err)
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)

    //ioutil.WriteFile("sina.json", body, 0666)
    //c.Infof(string(body))
    fmt.Fprintf(w, "%v", string(body))
    return
}

func CallbackSina(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    RequestSecret = r.FormValue("oauth_token")
    Code = r.FormValue("oauth_verifier")
    c.Infof("%v, %v", RequestSecret, Code)
    aTokenSina(w, r)
    return
}
