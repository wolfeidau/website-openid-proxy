package cookie

import (
	"net/http"
	"time"
)

// Config configures http.Cookie creation.
type Config struct {
	// Name is the desired cookie name.
	Name string
	// Domain sets the cookie domain. Defaults to the host name of the responding
	// server when left zero valued.
	Domain string
	// Path sets the cookie path. Defaults to the path of the URL responding to
	// the request when left zero valued.
	Path string
	// MaxAge=0 means no 'Max-Age' attribute should be set.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	// Cookie 'Expires' will be set (or left unset) according to MaxAge
	MaxAge int
	// HTTPOnly indicates whether the browser should prohibit a cookie from
	// being accessible via Javascript. Recommended true.
	HTTPOnly bool
	// Secure flag indicating to the browser that the cookie should only be
	// transmitted over a TLS HTTPS connection. Recommended true in production.
	Secure bool
}

// DefaultCookieConfig configures short-lived temporary http.Cookie creation.
var DefaultCookieConfig = Config{
	Name:     "login-temporary-cookie",
	Path:     "/",
	MaxAge:   60, // 60 seconds
	HTTPOnly: true,
	Secure:   true, // HTTPS only
}

// NewCookie returns a new http.Cookie with the given value and CookieConfig
// properties (name, max-age, etc.).
//
// The MaxAge field is used to determine whether an Expires field should be
// added for Internet Explorer compatibility and what its value should be.
func NewCookie(config Config, value string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     config.Name,
		Value:    value,
		Domain:   config.Domain,
		Path:     config.Path,
		MaxAge:   config.MaxAge,
		HttpOnly: config.HTTPOnly,
		Secure:   config.Secure,
	}
	// IE <9 does not understand MaxAge, set Expires if MaxAge is non-zero.
	if expires, ok := expiresTime(config.MaxAge); ok {
		cookie.Expires = expires
	}
	return cookie
}

// expiresTime converts a maxAge time in seconds to a time.Time in the future
// if the maxAge is positive or the beginning of the epoch if maxAge is
// negative. If maxAge is exactly 0, an empty time and false are returned
// (so the Cookie Expires field should not be set).
// http://golang.org/src/net/http/cookie.go?s=618:801#L23
func expiresTime(maxAge int) (time.Time, bool) {
	if maxAge > 0 {
		d := time.Duration(maxAge) * time.Second
		return time.Now().Add(d), true
	} else if maxAge < 0 {
		return time.Unix(1, 0), true // first second of the epoch
	}
	return time.Time{}, false
}
