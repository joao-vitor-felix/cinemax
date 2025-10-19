package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type contextKey string

const localizerKey contextKey = "localizer"

var locales = []string{"en", "pt-BR"}

func LocalizeMiddleware(bundle *i18n.Bundle) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			acceptLang := r.Header.Get("Accept-Language")
			if acceptLang == "" || !slices.Contains(locales, acceptLang) {
				acceptLang = "en"
			}
			if acceptLang != "" {
				langs := strings.Split(acceptLang, ",")
				if len(langs) > 0 {
					acceptLang = strings.Split(langs[0], ";")[0]
					if !slices.Contains(locales, acceptLang) {
						acceptLang = "en"
					}
				}
			}
			localizer := i18n.NewLocalizer(bundle, acceptLang)
			ctx := context.WithValue(r.Context(), localizerKey, localizer)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
	return f
}

func GetLocalizer(ctx context.Context) *i18n.Localizer {
	loc, ok := ctx.Value(localizerKey).(*i18n.Localizer)
	if !ok {
		return nil
	}
	return loc
}
