package auth

import (
    "path"
    "os"
    "io/ioutil"
    . "github.com/rugo/sacapi/modules/apilog"
    "crypto/sha256"
    "errors"
    "encoding/base64"
    "golang.org/x/oauth2"
    "encoding/json"
)

var (
    DataDir string = "/tmp/secrets" // TODO: Move to configuration
)

const (
    MIN_LEN_DEVICE_ID  = 5
    MIN_LEN_SECRET = 5
)



// Only authenticates if file [deviceId]/secret exists with content == [secret]
func AuthenticateByFile(deviceId, secret string) bool {
    if DeviceIdExists(deviceId) && CheckSecret(deviceId, secret) {
        Log.Info("Auth for device %s succeeded", deviceId)
        return true
    }
    Log.Info("Auth for device %s failed due to bad secret", deviceId)
    return false
}

func getDevicePath(deviceId string) string {
    return  path.Join(DataDir, deviceId)
}

func getDeviceSecretPath(deviceId string) string {
    return path.Join(getDevicePath(deviceId), "secret")
}

func getDeviceTokenPath(deviceId string) string {
    return path.Join(getDevicePath(deviceId), "token")
}

func DeviceIdExists(deviceId string) bool {
    if _, err := os.Stat(getDevicePath(deviceId)); os.IsNotExist(err) {
        Log.Info("Auth failed, deviceId %s not existent", deviceId)
        return false
    }
    return true
}

func getHashedDeviceSecret(deviceId string) (string, error) {
    if content, err := ioutil.ReadFile(getDeviceSecretPath(deviceId)); err == nil {
        return string(content[:len(content)-1]), nil
    }
    return "", errors.New("Secret not existent")
}

func CheckSecret(deviceId, secret string) bool {
    return CheckHashedSecret(deviceId, HashSecret(secret))
}

func CheckHashedSecret(deviceId, secretb64 string) bool {
    storedSecret, err := getHashedDeviceSecret(deviceId)
    if storedSecret == secretb64 && err == nil {
        return true
    }
    return false
}

func DeviceIsConnected(deviceId string) bool {
    if _, err := os.Stat(getDeviceTokenPath(deviceId)); os.IsNotExist(err) {
        return false
    }
    return true
}

func RegisterDevice(deviceId, secret string) error {
    if err := os.Mkdir(getDevicePath(deviceId), os.FileMode(0700)); err != nil {
        return err
    }
    if err := ioutil.WriteFile(getDeviceSecretPath(deviceId), []byte(HashSecret(secret)), 0600); err != nil {
        return err
    }
    return nil
}

func HashSecret(secret string) string {
    b := sha256.Sum256([]byte(secret))
    c := b[:]
    return base64.URLEncoding.EncodeToString(c)
}

func LoadToken(deviceId string) (*oauth2.Token, error) {
    f, err := os.Open(getDeviceTokenPath(deviceId))
    if err != nil {
        return nil, err
    }
    t := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(t)
    defer f.Close()
    return t, err
}

func SaveToken(deviceId string, token *oauth2.Token) error {
    f, err := os.Create(getDeviceTokenPath(deviceId))
    if err != nil {
        return err
    }
    defer f.Close()
    return json.NewEncoder(f).Encode(token)
}