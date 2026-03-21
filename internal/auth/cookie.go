package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
	// TODO:MAKE THESE ENV VARIABLES AND/OR MOVE IT TO THE CONFIG
	DefaultAccessDuration  = 15 * time.Minute
	DefaultRefreshDuration = 7 * 24 * time.Hour
)

func SetTokenCookies(w http.ResponseWriter, maker Maker, userID uuid.UUID, email string, isSecure bool) error {
	access, err := maker.Create(userID, email, AccessToken, DefaultAccessDuration)
	if err != nil {
		return err
	}
	refresh, err := maker.Create(userID, email, RefreshToken, DefaultRefreshDuration)
	if err != nil {
		return err
	}

	setCookie(w, AccessTokenCookie, access, DefaultAccessDuration, isSecure)
	setCookie(w, RefreshTokenCookie, refresh, DefaultRefreshDuration, isSecure)
	return nil
}

func ClearTokenCookies(w http.ResponseWriter) {
	setCookie(w, AccessTokenCookie, "", -time.Hour, false)
	setCookie(w, RefreshTokenCookie, "", -time.Hour, false)
}

func setCookie(w http.ResponseWriter, name, value string, dur time.Duration, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(dur.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}
