package strategy

import (
	"github.com/FGM/kurz/storage"
	"github.com/FGM/kurz/url"
	"testing"
)

// FIXME stop using constant credentials
func initTestStorage(t *testing.T) {
	var DSN = "root:@tcp(localhost:3306)/go_kurz_test"
	storage.Service.SetDSN(DSN)
	err := storage.Service.Open()
	if err != nil {
		t.Fatalf("Failed opening the test database: %+v", err)
	}
}

func TestBaseAlias(t *testing.T) {
	const BASE = "base"

	account := storage.User{
		DefaultStrategy: BASE,
	}
	strategy := Strategies[account.DefaultStrategy]
	if strategy.Name() != BASE {
		t.Fatalf("Strategy: expected %s, got %s", BASE, strategy.Name())
	}

	initTestStorage(t)
	// defers are executed LIFO
	defer storage.Service.Close()

	storage.Service.Truncate("shorturl")
	storage.Service.Truncate("longurl")
	storage.Service.AddToTruncateList("shorturl")
	storage.Service.AddToTruncateList("longurl")

	sourceUrl := url.LongUrl{
		Value: "http://www.example.com",
	}
	alias, err := strategy.Alias(&sourceUrl, storage.Service)
	if err != nil {
		t.Errorf("Failed during Alias(): %+v", err)
	}
	if alias.ShortFor.Id != sourceUrl.Id {
		t.Errorf("Aliasing does not point to proper long URL: expected %+v, got %+v", sourceUrl, alias.ShortFor)
	}

	if alias.Value != sourceUrl.Value {
		t.Errorf("Aliasing does not build the proper URL: expected %+v, got %+v", sourceUrl.Value, alias.Value)
	}
}

func TestUseCounts(t *testing.T) {
	const BASE = "base"

	account := storage.User{
		DefaultStrategy: BASE,
	}
	strategy := Strategies[account.DefaultStrategy]
	if strategy.Name() != BASE {
		t.Fatalf("Strategy: expected %s, got %s", BASE, strategy.Name())
	}

	initTestStorage(t)
	defer storage.Service.Close()
	storage.Service.Truncate("shorturl")
	storage.Service.Truncate("longurl")
	storage.Service.AddToTruncateList("shorturl")
	storage.Service.AddToTruncateList("longurl")

	initialCount := strategy.UseCount(storage.Service)
	if initialCount != 0 {
		t.Errorf("Found %d record(s) in test database, expecting none.", initialCount)
	}

	sourceUrl := url.LongUrl{
		Value: "http://www.example.com",
	}
	_, err := strategy.Alias(&sourceUrl, storage.Service)
	if err != nil {
		t.Errorf("Failed during Alias(): %+v", err)
	}

	nextCount := strategy.UseCount(storage.Service)
	if nextCount != initialCount+1 {
		t.Errorf("Found %d record(s) in test database, expecting %d.", nextCount, initialCount+1)
	}

	sourceUrl = url.LongUrl{
		Value: "http://www2.example.com",
	}
	_, err = strategy.Alias(&sourceUrl, storage.Service)
	if err != nil {
		t.Errorf("Failed during Alias(): %+v", err)
	}

	nextCount = strategy.UseCount(storage.Service)
	if nextCount != initialCount+2 {
		t.Errorf("Found %d record(s) in test database, expecting %d.", nextCount, initialCount+2)
	}
}
