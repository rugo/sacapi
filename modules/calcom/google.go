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
	"github.com/rugo/sacapi/modules/maps"
	"time"
)

var (
	oauthConfig *oauth2.Config
)

const (
	MAX_MAP_DURATION_MINS = 8 * 60  // Max offset used, in case appointment has location
	MAP_ORIGIN = "Konstanz, Deutschland" // This is just for testing, will be replaced by geolocation stuff
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
	mapsDone := make(chan bool, 1)  // for signaling when maps api call is done
	var (
		bReminderMin int64 = 0 // the offset to the appointment, calculated by reminders and distance
		errMaps error = nil
		durationMapsMins int64 = 0
	)
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
		// Skip all day events
		for _, i := range events.Items {
			if entry.Start.DateTime != "" {
				break;
			}
			entry = i
		}

		if entry.Start.DateTime == "" {
			return ctx, ErrNoAppointments
		}
		if entry.Location != "" {
			go func() {
				durationMapsMins, errMaps = maps.GetDurationMins(ctx, MAP_ORIGIN, entry.Location)
				mapsDone <- true
			}()
		} else {
			mapsDone <- true
		}

		startTime, err := time.Parse(time.RFC3339, entry.Start.DateTime)
        /* Get reminder times */
        reminders := entry.Reminders.Overrides
        // Default reminders (in case used)
        if entry.Reminders.UseDefault {
            calendarList, errList := srv.CalendarList.Get("primary").Do()
            if errList != nil {
                Log.Error("Unable to get CalendarList entry of primary")
                return ctx, ErrCommunicationError
            }
            reminders = calendarList.DefaultReminders
        }

        for _, reminder := range reminders {
            if reminder.Minutes > bReminderMin {
                bReminderMin = reminder.Minutes
            }
        }

		if err != nil {
			Log.Error("Could not parse time %s", entry.Start.DateTime)
			return ctx, ErrCommunicationError
		}

		// Wait for Maps api answer
		select {
		case <-ctx.Done():
			return ctx, ctx.Err()
		case <- mapsDone:
			if errMaps == nil && durationMapsMins < MAX_MAP_DURATION_MINS {
				bReminderMin = bReminderMin + durationMapsMins
			}
		}

		ctx = NewContext(ctx,  data.ClockInfoPackage{
			Appointment: data.Appointment{
				Time: startTime.Unix() - bReminderMin * 60,
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
