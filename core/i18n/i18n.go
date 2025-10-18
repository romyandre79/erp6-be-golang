package i18n

import (
	"encoding/json"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load semua file locale
	files := []string{
		"./locales/en.json",
		"./locales/id.json",
	}

	for _, f := range files {
		if _, err := Bundle.LoadMessageFile(f); err != nil {
			log.Fatalf("Failed to load locale file %s: %v", f, err)
		}
	}
}

// Helper untuk translate
func Translate(lang string, messageID string, data map[string]interface{}) string {
	localizer := i18n.NewLocalizer(Bundle, lang)
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	return msg
}
