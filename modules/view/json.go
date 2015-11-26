package view
import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/rugo/sacapi/modules/data"
    "strconv"
    "golang.org/x/net/context"
    "github.com/rugo/sacapi/modules/calcom"
)

func GetJSONMessage(w rest.ResponseWriter, r *rest.Request) {
    time, err := strconv.Atoi(r.PathParam("time"))
    if err != nil {
        rest.Error(w, "Time has to be a number", 400)
        return
    }
    t := data.ClockInfo{
        Appointment: data.Appointment{
            Time: time,
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
    nextEntry, err := calcom.GetNextGoogleCalendarEntry(ctx, deviceId)

    if err != nil {
        rest.Error(w, "Could not read calendar entries", 500)
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteJson(nextEntry)
}