package auth

import (
    "path"
    "os"
    "io/ioutil"
    . "github.com/rugo/sacapi/modules/apilog"
)

var (
    DataDir string = "/tmp/secrets" // TODO: Move to configuration
)


// Only authenticates if file [deviceId]/secret exists with content == [secret]
func AuthenticateByFile(deviceId string, secret string) bool {
    pathToSecret := path.Join(DataDir, deviceId, "secret")
    if _, err := os.Stat(pathToSecret); os.IsNotExist(err) {
        Log.Info("Auth failed, path %s not existent", pathToSecret)
        return false
    }

    if content, err := ioutil.ReadFile(pathToSecret); err == nil {
        if string(content[:len(content)-1]) == secret {
            Log.Info("Auth for id %s successful", deviceId)
            return true
        }
    }
    Log.Info("Auth for device %s failed due to bad secret", deviceId)
    return false
}
