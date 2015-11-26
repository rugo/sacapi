package view
import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/rugo/sacapi/modules/data"
    "strconv"
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

/*

*/
func RegisterDevice(w rest.ResponseWriter, r *rest.Request) {
    content := &RegisterDeviceRequest{}
    err := r.DecodeJsonPayload(&content)
    if err != nil {
        rest.Error(w, "Invalid request", 400)
        return
    }
    context := context.Background()
    calendar.Register(ctx, content)
}