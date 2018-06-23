package jp

import (
	"context"
	"fmt"
	"log"
	"testing"

	"cloud.google.com/go/translate"
	xlang "golang.org/x/text/language"
)

func TestJapaneseTranslate(t *testing.T) {
	const text = "安倍首相は３０日、野党党首との党首討論で、トランプ米政権が検討している輸入車への高関税措置について、「同盟国の日本に高関税を課すのは極めて理解しがたく、受け入れることはできない」と述べ、反対する意向を明言した。"
	ctx := context.Background()
	tc, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create translate client: %s", err)
	}
	ts, err := tc.Translate(ctx, []string{text}, xlang.English, nil)
	if err != nil {
		log.Fatalf("failed to translate text: %s", err)
	}
	if len(ts) != 1 {
		log.Fatalf("unexpected count. want 1, got %d", len(ts))
	}
	fmt.Printf("Translation: %s\n", ts[0].Text)
}
