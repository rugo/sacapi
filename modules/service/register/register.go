package main

import (
    . "github.com/rugo/sacapi/modules/apilog"
    "github.com/rugo/sacapi/modules/auth"
    "net/http"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/calendar/v3"
    "io/ioutil"
)

var (
    oauthConfig *oauth2.Config
)



func respondBadRequest(w http.ResponseWriter, reason string) {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(reason))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    deviceId := r.PostFormValue("deviceid")
    secret := r.PostFormValue("secret")

    if len(deviceId) < auth.MIN_LEN_DEVICE_ID || len(secret) < auth.MIN_LEN_SECRET {
        respondBadRequest(w, "Arguments have to wrong format")
        return
    }

    Log.Critical("Received register request from device %s", deviceId)

    if auth.DeviceIdExists(deviceId) && auth.DeviceIsConnected(deviceId){
        respondBadRequest(w, "Device already registered")
        Log.Info("Already registered device %s entered /register", r.Form["deviceId"])
        return
    }

    err := auth.RegisterDevice(deviceId, secret)
    if err != nil {
        Log.Error(err.Error())
    }
    http.SetCookie(w, &http.Cookie{Name: "deviceid", Value: deviceId})
    http.SetCookie(w, &http.Cookie{Name: "secret", Value: auth.HashSecret(secret)})
    http.Redirect(w, r, oauthConfig.AuthCodeURL(""), http.StatusFound)
}

func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) {
    deviceIdC, errDev := r.Cookie("deviceid")
    secretb64C, errSec := r.Cookie("secret")

    if errDev != nil || errSec != nil {
        respondBadRequest(w, "You need to register first.")
    }

    deviceId := deviceIdC.Value
    secretb64 := secretb64C.Value

    if auth.DeviceIdExists(deviceId) && auth.CheckHashedSecret(deviceId, secretb64) {
        Log.Info("Received callback for device %s", deviceId)
        //Get the code from the response
        code := r.FormValue("code")

        // Exchange for token
        token, err := oauthConfig.Exchange(oauth2.NoContext, code)

        if err != nil {
            Log.Error("Could not exchange code to token for device %s", )
        }

        err = auth.SaveToken(deviceId, token)
        if err != nil {
            Log.Error("Could not save token for device %s", deviceId)
        }
    }
}

func main() {
    b, err := ioutil.ReadFile("/tmp/google_api_secret.json")
    if err != nil {
        Log.Fatalf("Unable to read client secret file: %v", err)
    }

    oauthConfig, err = google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
    if err != nil {
        Log.Fatalf("Unable to parse client secret file to config: %v", err)
    }

    http.HandleFunc("/register", registerHandler)
    http.HandleFunc("/oauth2callback", oauthCallbackHandler)
    http.ListenAndServe(":8080", nil)
}