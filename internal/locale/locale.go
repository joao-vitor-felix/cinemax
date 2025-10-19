package locale

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewBundle() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("internal/locale/en.json")
	bundle.LoadMessageFile("internal/locale/pt-BR.json")
	return bundle
}

func GetMessage(localizer *i18n.Localizer, messageID string) string {
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	return msg
}

func GetTemplateMessage(localizer *i18n.Localizer, messageID string, template map[string]any) string {
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: template,
	})
	return msg
}
