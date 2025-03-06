package translation

import (
	"context"
	"fmt"
	"google.golang.org/api/option"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func TranslateText(targetLanguage, text string) (string, error) {
	// text := "The Go Gopher is cute"
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %w", err)
	}

	client, err := translate.NewClient(ctx, option.WithAPIKey("AIzaSyBZ3zvvD8dF-q3pdZX4hqVmAfLUpf4imGI"))
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("Translate: %w", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}
