package locales

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("internal/locales/en.json")
	bundle.LoadMessageFile("internal/locales/pt-BR.json")
	return bundle
}

func NewLocalizer(bundle *i18n.Bundle, langs ...string) *i18n.Localizer {
	localizer := i18n.NewLocalizer(bundle, langs...)
	return localizer
}
