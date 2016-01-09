package maps

import (
	. "github.com/rugo/sacapi/modules/apilog"
	"googlemaps.github.io/maps"
	"golang.org/x/net/context"
	"time"
	"io/ioutil"
)

const (
	API_KEY_FILE = "/tmp/maps_api_key"
)
var (
	apiKey = ""
)

func initApi() error {
	if apiKey == "" {
		key, err := ioutil.ReadFile(API_KEY_FILE)
		if err != nil {
			Log.Error("Could not read MAPS API Key file")
			return err
		}
		apiKey = string(key)
	}
	return nil
}

func getDuration(ctx context.Context, origin, destination string) (time.Duration, error) {
	err := initApi()
	if err != nil {
		return time.Duration(0), err
	}

	c, err := maps.NewClient(maps.WithAPIKey(apiKey))

	if err != nil {
		Log.Error("Could not create MAPS client %s", err)
		return time.Duration(0), err
	}

	r := &maps.DirectionsRequest{
		Origin: origin,
		Destination: destination,
	}
	resp, _, err := c.Directions(context.Background(), r)

	if err != nil || len(resp) == 0 {
		Log.Error("Could not get directions for route %s to %s, err:%s", origin, destination, err)
		return time.Duration(0), err
	}

	return resp[0].Legs[0].Duration, nil
}

func GetDurationMins(ctx context.Context, origin, destination string) (int64, error) {
	dur, err := getDuration(ctx, origin, destination)
	if err != nil {
		return 0, err
	}
	return int64(dur.Minutes()), nil
}