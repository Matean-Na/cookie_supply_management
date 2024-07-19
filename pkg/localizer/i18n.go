package localizer

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"strings"
)

type LocalizeError struct {
	Source  error
	Message string
	Data    interface{}
}

func Localize(ctx *gin.Context, messageID string, data interface{}) string {
	l, exists := ctx.Get("localizer")
	if !exists {
		return "error"
	}

	localize := l.(*i18n.Localizer)
	resultMessage, err := localize.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		defaultMessage, err := localize.Localize(&i18n.LocalizeConfig{
			MessageID:    "exception:default-message",
			TemplateData: data,
		})
		if err != nil {
			return "error"
		}
		return defaultMessage
	}
	return resultMessage
}

func (e LocalizeError) Error() string {
	if e.Source == nil {
		return e.Message
	}
	return e.Source.Error()
}

func NewLocalizeError(err error, message string, data interface{}) *LocalizeError {
	return &LocalizeError{
		Source:  err,
		Message: message,
		Data:    data,
	}
}

func (e LocalizeError) Localize(ctx *gin.Context) string {
	return Localize(ctx, e.Message, e.Data)
}

func DefaultError(ctx *gin.Context) string {
	return Localize(ctx, "exception:default-message", nil)
}

func GetLanguageFromHeader(input string) string {
	var languageStr string = strings.Split(input, ";")[0]
	var languageFirstPart = strings.Split(languageStr, ",")

	if len(languageFirstPart[0]) <= 0 {
		return "en"
	}

	languageStr = languageFirstPart[len(languageFirstPart)-1]

	return strings.ToLower(languageStr)
}

func Init() error {
	locales := []string{"ru_RU", "en_US", "ky_KG"}

	bundle := i18n.NewBundle(language.Russian)

	for _, locale := range locales {
		filePath := fmt.Sprintf("translations/%s.json", locale)

		_, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
		_, err = bundle.LoadMessageFile(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}
