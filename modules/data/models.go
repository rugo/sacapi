package data

type Appointment struct {
    Time int64 `json:"time"`
    Name string `json:"name"`
    Description string `json:"description"`
    Location string `json:"location"`
}

type ClockInfoPackage struct {
    Appointment Appointment `json:"appointment"`
    Timezone string `json:"timezone"`
    Apivers int `json:"apivers"`
}