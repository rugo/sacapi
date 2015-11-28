package calcom
import (
	"golang.org/x/net/context"
	"github.com/rugo/sacapi/modules/data"
	"errors"
)

/*
Calendar communication
 */

var (
	ErrCommunicationError = errors.New("Communication error with API")
	ErrNoAppointments = errors.New("The user has no appointments")
)

type key int
const (
	appointmentKey key = 0
)

func NewContext(ctx context.Context, app data.ClockInfoPackage) context.Context {
	return context.WithValue(ctx, appointmentKey, app)
}

func FromContext(ctx context.Context) (data.ClockInfoPackage, bool) {
	app, ok := ctx.Value(appointmentKey).(data.ClockInfoPackage)
	return app, ok
}

func InitApis() {
	initGoogleCalendarApi() /* TODO: move API methods to dict */
}