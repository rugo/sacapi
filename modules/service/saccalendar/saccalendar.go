package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/rugo/sacapi/modules/auth"
    "github.com/rugo/sacapi/modules/view"
    "github.com/rugo/sacapi/modules/apilog"
    "log"
    "net/http"
    "github.com/rugo/sacapi/modules/calcom"
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
            Authenticator: auth.AuthenticateByFile,
            Authorizator: auth.AuthorizeRequest})
    router, err := rest.MakeRouter(
        rest.Get("/test/:time", view.GetJSONMessage),
        rest.Get("/calendar/next/:id", view.GetNextCalendarEntry),
    )
    if err != nil {
        log.Fatal(err)
    }
    api.SetApp(router)
    apilog.Init()
    calcom.InitApis()
    log.Fatal(
        http.ListenAndServeTLS(":1443",
            "/etc/sac/keys/cert.pem",
            "/etc/sac/keys/key.pem",
            api.MakeHandler(),
        ),
    )
}