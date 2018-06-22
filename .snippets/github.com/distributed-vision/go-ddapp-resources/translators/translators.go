package translators

import (
	"context"
	"fmt"
	"sync"

	"github.com/distributed-vision/go-resources/ids"
)

type TranslationFunction func(translationContext context.Context, fromId ids.Identifier, fromValue interface{}) (chan interface{}, chan error)

type translationEntry struct {
	fromType   ids.TypeIdentifier
	toType     ids.TypeIdentifier
	translator TranslationFunction
}

var translators = make(map[string]map[string]*translationEntry)
var translatorMutex = sync.Mutex{}

func Translate(translationContext context.Context, fromType ids.TypeIdentifier, fromId ids.Identifier, fromValue interface{}, toType ids.TypeIdentifier) (chan interface{}, chan error) {
	var translator TranslationFunction
	//fmt.Printf("TRNS: %v -> %v: %v\n", fromType, toType, fromValue)

	translatorMutex.Lock()
	if entryMap, ok := translators[string(fromType.Value())]; ok {
		if entry, ok := entryMap[string(toType.Value())]; ok {
			translator = entry.translator
		}
	}
	translatorMutex.Unlock()

	if translator != nil {
		//fmt.Printf("TRANS\n")
		return translator(translationContext, fromId, fromValue)
	}
	//fmt.Printf("NO TRANS\n")
	cres := make(chan interface{}, 1)
	cerr := make(chan error, 1)

	cerr <- fmt.Errorf("Can't find transtaor for: %v to %v", fromType, toType)
	close(cres)
	close(cerr)
	return cres, cerr
}

func Register(translationContext context.Context, fromType ids.TypeIdentifier, toType ids.TypeIdentifier, translator TranslationFunction) TranslationFunction {
	var previous TranslationFunction
	//fmt.Printf("TRNSREG: %v -> %v\n", fromType, toType)
	translatorMutex.Lock()
	if entryMap, ok := translators[string(fromType.Value())]; ok {
		if entry, ok := entryMap[string(toType.Value())]; ok {
			previous = entry.translator
			entry.translator = translator
		} else {
			entryMap[string(toType.Value())] = &translationEntry{fromType, toType, translator}
		}
	} else {
		entryMap = make(map[string]*translationEntry)
		entryMap[string(toType.Value())] = &translationEntry{fromType, toType, translator}
		translators[string(fromType.Value())] = entryMap
	}
	translatorMutex.Unlock()

	return previous
}
