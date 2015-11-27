package calcom
import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	. "github.com/rugo/sacapi/modules/apilog"
	"github.com/rugo/sacapi/modules/auth"
	"github.com/rugo/sacapi/modules/data"
	"time"
	"errors"
)

var (
	oauthConfig *oauth2.Config
)

func GetNextGoogleCalendarEntry(ctx context.Context, deviceId string) (data.ClockInfo, error) {
	/* ToDo: remove test code */
	// Not to be done in every request, JUST FOR TESTING!!
	b, err := ioutil.ReadFile("/etc/sac/google_api_secret.json")
	if err != nil {
		Log.Fatalf("Unable to read client secret file: %v", err)
	}

	oauthConfig, err = google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		Log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	token, err := auth.LoadToken(deviceId)
	if err != nil {
		Log.Error("Could not load token for device %s", deviceId)
		return data.ClockInfo{}, err
	}

	client := oauthConfig.Client(oauth2.NoContext, token)
	srv, err := calendar.New(client)

	if err != nil {
		Log.Error("Unable to retrieve calendar Client %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
	SingleEvents(true).TimeMin(t).MaxResults(1).OrderBy("startTime").Do()
	if err != nil {
		Log.Error("Unable to retrieve next ten of the user's events. %v", err)
		return data.ClockInfo{}, err
	}

	if len(events.Items) > 0 {
		entry := events.Items[0]
		startTime, err := time.Parse(time.RFC3339, entry.Start.DateTime)

		if err != nil {
			Log.Error("Could not parse time %s", entry.Start.DateTime)
			return data.ClockInfo{}, errors.New("Communication error with API")
		}

		nextEntry := data.ClockInfo{
			Appointment: data.Appointment{
				Time: startTime.Unix(),
				Name: entry.Summary,
				Description: entry.Description,
			},
			Timezone: events.TimeZone,
			Apivers: 0,
		}
		return nextEntry, nil
	} else {
		errMsg := "Device %s user has no appointments."
		Log.Error(errMsg, deviceId)
		return data.ClockInfo{}, errors.New(errMsg)
	}

}
