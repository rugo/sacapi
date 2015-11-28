package view
import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/rugo/sacapi/modules/data"
    "strconv"
    "golang.org/x/net/context"
    "github.com/rugo/sacapi/modules/calcom"
    . "github.com/rugo/sacapi/modules/apilog"
)

const (
    CONTENT_TYPE = "application/json; charset=utf-8"
    ERROR_CODE = 500 /* always return 500 if smth goes wrong here */
)
func GetJSONMessage(w rest.ResponseWriter, r *rest.Request) {
    time, err := strconv.Atoi(r.PathParam("time"))
    if err != nil {
        rest.Error(w, "Time has to be a number", 400)
        return
    }
    t := data.ClockInfoPackage{
        Appointment: data.Appointment{
            Time: int64(time),
            Name: "Meeting",
            Description: "Nöpe! Chuck Testa!",
        },
        Timezone: "UTC+01:00",
        Apivers: 0,
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteJson(t)
}

func GetNextCalendarEntry(w rest.ResponseWriter, r *rest.Request) {
    ctx := context.Background()
    deviceId := r.PathParam("id")
    ctx, err := calcom.GetNextGoogleCalendarEntry(ctx, deviceId)

    if err != nil {
        Log.Error(err.Error())
        rest.Error(w, "Could not read calendar entries", ERROR_CODE)
        return
    }

    answer, ok := calcom.FromContext(ctx)

    if !ok {
        /* should never be reached */
        rest.Error(w, "Could not access next calendar entry", ERROR_CODE)
    }

    w.Header().Set("Content-Type", CONTENT_TYPE)
    w.WriteJson(answer)
}