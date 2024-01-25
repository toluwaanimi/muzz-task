package constants

import "time"

type MiddlewareContextKey string

func (m MiddlewareContextKey) String() string {
	return string(m)
}

var (
	AuthenticatedAccountContextKey      MiddlewareContextKey = "account"
	AuthenticatedSessionTokenContextKey MiddlewareContextKey = "token"
)

const DefaultUserCount = 100

const DefaultPassword = "password"

const UserCollection = "users"
const MatchCollection = "matches"
const SwipeCollection = "swipes"

const (
	TokenIssuer         = "Muzz Dating"
	TokenSubject        = "Muzz Dating Token"
	TokenAudience       = "https://muzz.com"
	TokenJWTID          = "Muzz Dating"
	DefaultTokenExpires = 24 * time.Hour
)

type LondonCoordinates struct {
	MinLat float64
	MaxLat float64
	MinLon float64
	MaxLon float64
}

var NorthLondonCoordinates = LondonCoordinates{
	MinLat: 51.5244,
	MaxLat: 51.6722,
	MinLon: -0.2076,
	MaxLon: 0.1698,
}

var WhitelistedRoutes = map[string]struct{}{
	"/":                  {},
	"/docs":              {},
	"developers":         {},
	"/docs/swagger.yaml": {},
	"/user/create":       {},
	"/login":             {},
}
