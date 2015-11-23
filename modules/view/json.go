package view
import (
    "github.com/ant0ine/go-json-rest/rest"
    "github.com/rugo/sacapi/modules/data"
)

func GetJSONMessage(w rest.ResponseWriter, r *rest.Request) {
    t := data.ClockInfo{
        Appointment: data.Appointment{
            Time: 123456,
            Name: "Meeting",
            Description: "Nope! Chuck Testa!",
        },
        Timezone: "UTC+01:00",
        Apivers: 0,
    }
    w.WriteJson(t)
}