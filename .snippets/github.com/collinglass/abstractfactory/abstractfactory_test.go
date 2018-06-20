package abstractfactory

import (
	"testing"
)

func TestCreateStore(t *testing.T) {
	var factory FancyFactory
	store := ErectStore(factory)

	_, err := store.GetStock(1)
	if err != nil {
		t.Errorf("TestCreateStore failed: %s \n", err)
	}
}

func TestCreateMetalProduct(t *testing.T) {
	var factory FancyFactory
	mp := factory.CreateMetalProduct()

	if mp == nil {
		t.Errorf("TestCreateMetalProduct failed: %s \n", "didn't create anything")
	}
}

func TestCreatePaperProduct(t *testing.T) {
	var factory FancyFactory
	mp := factory.CreatePaperProduct()

	if mp == nil {
		t.Errorf("TestCreatePaperProduct failed: %s \n", "didn't create anything")
	}
}
