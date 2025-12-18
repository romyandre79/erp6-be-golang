package i18n

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var Bundle *i18n.Bundle

func Init() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

    // Debug CWD
    if cwd, err := os.Getwd(); err == nil {
        log.Printf("Start I18N Check...")
        log.Printf("Current Working Directory: %s", cwd)
        
        files, _ := os.ReadDir("./locales")
        for _, f := range files {
            log.Printf("Found in locales: %s", f.Name())
        }
    }

	// Load semua file locale
	files := []string{
		"locales/en.json",
		"locales/id.json",
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
	if msg == "" {
		msg = messageID
	}
	return msg
}
