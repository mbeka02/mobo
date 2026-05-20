package auth

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
)

func SetTokenCookies(w http.ResponseWriter, maker Maker, userID uuid.UUID, email, role string, isSecure bool, accessDuration, refreshDuration time.Duration) error {
	access, err := maker.Create(userID, email, role, AccessToken, accessDuration)
	if err != nil {
		return err
	}
	refresh, err := maker.Create(userID, email, role, RefreshToken, refreshDuration)
	if err != nil {
		return err
	}

	setCookie(w, AccessTokenCookie, access, accessDuration, isSecure)
	setCookie(w, RefreshTokenCookie, refresh, refreshDuration, isSecure)
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
