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
)

var (
	oauthConfig *oauth2.Config
)

func initGoogleCalendarApi() {
	b, err := ioutil.ReadFile("/etc/sac/google_api_secret.json")
	if err != nil {
		Log.Fatalf("Unable to read client secret file: %v", err)
	}

	oauthConfig, err = google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		Log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	Log.Info("Read Google API oauth config.")
}

func GetNextGoogleCalendarEntry(ctx context.Context, deviceId string) (context.Context, error) {
	token, err := auth.LoadToken(deviceId)
	if err != nil {
		Log.Error("Could not load token for device %s", deviceId)
		return ctx, err
	}

	Log.Info("Using oauth config %s", oauthConfig.ClientID)

	client := oauthConfig.Client(oauth2.NoContext, token)
	srv, err := calendar.New(client)

	if err != nil {
		Log.Error("Unable to retrieve calendar Client %v", err)
		return ctx, err
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
	SingleEvents(true).TimeMin(t).MaxResults(1).OrderBy("startTime").Do()
	if err != nil {
		Log.Error("Unable to retrieve user event. %v", err)
		return ctx, err
	}

	if len(events.Items) > 0 {
		entry := events.Items[0]
		startTime, err := time.Parse(time.RFC3339, entry.Start.DateTime)

		if err != nil {
			Log.Error("Could not parse time %s", entry.Start.DateTime)
			return ctx, ErrCommunicationError
		}

		ctx = NewContext(ctx,  data.ClockInfoPackage{
			Appointment: data.Appointment{
				Time: startTime.Unix(),
				Name: entry.Summary,
				Description: entry.Description,
				Location: entry.Location,
			},
			Timezone: events.TimeZone,
			Apivers: 0,
		})
		return ctx, nil
	} else {
		return ctx, ErrNoAppointments
	}
}
