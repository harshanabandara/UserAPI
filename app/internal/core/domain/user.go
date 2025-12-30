package domain

import "encoding/json"

type UserStatus int

const (
	_ UserStatus = iota // We don't want this one.
	ACTIVE
	INACTIVE
)

var statusMap = map[UserStatus]string{
	ACTIVE:   "active",
	INACTIVE: "inactive",
}

var stringToStatusMap = map[string]UserStatus{
	"active":   ACTIVE,
	"inactive": INACTIVE,
}

func (u *UserStatus) UnmarshalJSON(bytes []byte) error {
	//Convert "status" : "active" to UserStatus
	var str string
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}
	*u = stringToStatusMap[str]
	return nil
}

func (u *UserStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *UserStatus) String() string {
	return statusMap[*u]
}

type User struct {
	UserID    string     `json:"userId,omitempty"`
	FirstName string     `json:"firstName,omitempty"`
	LastName  string     `json:"lastName,omitempty"`
	Email     string     `json:"email,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Age       int        `json:"age,omitempty"`
	Status    UserStatus `json:"status,omitempty"`
}
