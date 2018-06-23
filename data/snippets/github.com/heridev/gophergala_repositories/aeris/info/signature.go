package info

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func (i *Info) DecryptSignatures() error {

	if i.decryptedSignatures {
		return nil
	}

	res, err := http.Get(i.playerJsUrl)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	// body contains javascript code containing decryption info, we're about to parse those
	// and use them to decrypt signatures
	decryption, err := extractDecryption(body)
	if err != nil {
		return err
	}

	for _, stream := range i.streams {

		stream.signature = decryption.run(stream.signature)

		err = stream.buildSignatureUrl(stream.signature)
		if err != nil {
			return err
		}
	}

	i.decryptedSignatures = true

	return nil
}

func extractDecryption(js []byte) (chain, error) {

	actionString, object, err := extractMethodCalls(js)
	if err != nil {
		return nil, err
	}

	methodMapping, err := extractMethodMapping(object, js)
	if err != nil {
		return nil, err
	}

	actionCalls := strings.Split(actionString, ";")

	decryption, err := buildDecryptionChain(methodMapping, actionCalls)
	if err != nil {
		return decryption, err
	}

	return decryption, nil
}

type action struct {
	method method
	param  int
}

type chain []*action

func (c *chain) run(sig string) string {

	for _, action := range *c {
		sig = action.method.handler(sig, action.param)
	}

	return sig
}

var methodCallInfoRegex = regexp.MustCompile(`(?i)[a-z0-9]{2}\.([a-z0-9]{2})\([a-z],([0-9]+)\)`)

func buildDecryptionChain(methods map[string]method, jsCalls []string) (chain, error) {

	var actions chain

	for _, call := range jsCalls {

		match := methodCallInfoRegex.FindStringSubmatch(call)
		if match == nil {
			return nil, errors.New("Could not match info extraction regex against method call")
		}

		name := match[1]
		param, _ := strconv.Atoi(match[2])

		actions = append(actions, &action{
			method: methods[name],
			param:  param,
		})
	}

	return actions, nil
}

var methodsRegex = map[string]string{
	"reverse": `([a-z0-9]{2}):function\([a-z]\)\{[a-z]\.reverse\(\)\}`,
	"swap":    `([a-z0-9]{2}):function\([a-z],[a-z]\)\{var [a-z]=[a-z]\[[0-9]\];[a-z]\[[0-9]\]=[a-z]\[[a-z]%[a-z]\.length\];[a-z]\[[a-z]\]=[a-z]\}`,
	"splice":  `([a-z0-9]{2}):function\([a-z],[a-z]\)\{[a-z]\.splice\([0-9],[a-z]\)\}`,
}

func objectMethodExtractRegex(objectName string) (*regexp.Regexp, error) {

	var methodArray []string
	for _, regexStr := range methodsRegex {
		methodArray = append(methodArray, regexStr)
	}

	regex, err := regexp.Compile(`(?i)var ` + objectName + `=\{((?:(?:` + strings.Join(methodArray, "|") + `)(?:,|\})?)+)`)

	return regex, err
}

type method struct {
	name       string
	definition string
	handler    handler
}

type handler func(in string, param int) string

var methodsRegexToHandler = map[*regexp.Regexp]handler{
	regexp.MustCompile(`(?i)` + methodsRegex["reverse"]): reverseHandler,
	regexp.MustCompile(`(?i)` + methodsRegex["swap"]):    swapHandler,
	regexp.MustCompile(`(?i)` + methodsRegex["splice"]):  spliceHandler,
}

func extractMethodMapping(object string, js []byte) (map[string]method, error) {

	regex, err := objectMethodExtractRegex(object)
	if err != nil {
		return nil, err
	}

	match := regex.FindSubmatch(js)
	if match == nil {
		return nil, errors.New("Couldn't match object method extraction regex against js body")
	}

	definitions := string(match[1])

	methods := make(map[string]method)

	for regex, handler := range methodsRegexToHandler {
		match := regex.FindStringSubmatchIndex(definitions)
		if match != nil {

			definition := definitions[match[0]:match[1]]

			name := definitions[match[0] : match[0]+2]

			methods[name] = method{
				name:       name,
				definition: definition,
				handler:    handler,
			}
		}
	}

	return methods, nil
}

var actionsExtractRegex = regexp.MustCompile(`(?i)[a-z]=[a-z]\.split\(""\);((?:([a-z]{2})\.[a-z0-9]{2}\([a-z],[0-9]+\);)+)return [a-z]\.join\(""\)`)

func extractMethodCalls(js []byte) (string, string, error) {
	match := actionsExtractRegex.FindSubmatch(js)
	if match == nil {
		return "", "", errors.New("Could not match action extraction regex against js body")
	}

	actionString := string(match[1])
	object := string(match[2])

	actionString = strings.Trim(actionString, ";")

	return actionString, object, nil
}

func reverseHandler(sig string, _ int) string {

	runes := []rune(sig)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func swapHandler(sig string, pos int) string {
	runes := []rune(sig)

	temp := runes[0]

	runes[0] = runes[pos%len(runes)]

	runes[pos] = temp

	return string(runes)
}

func spliceHandler(sig string, pos int) string {
	runes := []rune(sig)
	runes = runes[pos:]
	return string(runes)
}
