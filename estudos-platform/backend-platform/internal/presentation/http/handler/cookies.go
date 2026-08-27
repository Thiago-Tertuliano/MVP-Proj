package handler

import "net/http"

const (
	cookieAccess  = "access_token"
	cookieRefresh = "refresh_token"
)

type AuthCookies struct {
	AccessMaxAge  int
	RefreshMaxAge int
	Secure        bool
}

func (c AuthCookies) withDefaults() AuthCookies {
	if c.AccessMaxAge <= 0 {
		c.AccessMaxAge = 15 * 60
	}
	if c.RefreshMaxAge <= 0 {
		c.RefreshMaxAge = 168 * 3600
	}
	return c
}

func (c AuthCookies) gravar(w http.ResponseWriter, access, refresh string) {
	c = c.withDefaults()
	http.SetCookie(w, &http.Cookie{
		Name:     cookieAccess,
		Value:    access,
		Path:     "/",
		MaxAge:   c.AccessMaxAge,
		HttpOnly: true,
		Secure:   c.Secure,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     cookieRefresh,
		Value:    refresh,
		Path:     "/api/v1/auth",
		MaxAge:   c.RefreshMaxAge,
		HttpOnly: true,
		Secure:   c.Secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func (c AuthCookies) limpar(w http.ResponseWriter) {
	expired := func(name, path string) *http.Cookie {
		return &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     path,
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   c.Secure,
			SameSite: http.SameSiteLaxMode,
		}
	}
	http.SetCookie(w, expired(cookieAccess, "/"))
	http.SetCookie(w, expired(cookieRefresh, "/api/v1/auth"))
}

func cookieValue(r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return c.Value
}
