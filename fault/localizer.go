package fault

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var Localization *i18n.Bundle

func InitLocales(locales ...string) error {
	Localization = i18n.NewBundle(language.Russian)
	Localization.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	for _, lang := range locales {
		if _, err := Localization.LoadMessageFile(lang); err != nil {
			return err
		}
	}
	return nil
}
