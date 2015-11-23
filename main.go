package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/rugo/sacapi/modules/auth"
    "github.com/rugo/sacapi/modules/view"
    "github.com/rugo/sacapi/modules/apilog"
    "log"
    "net/http"
)


/*
Microservice for accessing data
*/
func main() {
    api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)
    api.Use(
        &rest.AuthBasicMiddleware{
            Realm: "Smart Alarm Clock REST API",
            Authenticator: auth.AuthenticateByFile})
    router, err := rest.MakeRouter(
        rest.Get("/test", view.GetJSONMessage),
    )
    if err != nil {
        log.Fatal(err)
    }
    api.SetApp(router)
    apilog.Init()
    log.Fatal(http.ListenAndServeTLS(":1443", "keys/cert.pem", "keys/key.pem", api.MakeHandler()))
}