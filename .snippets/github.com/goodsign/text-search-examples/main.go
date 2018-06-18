package main

import (
    "fmt"
    "io/ioutil"
    "regexp"
    "github.com/goodsign/libtextcat"
    "github.com/goodsign/snowball"
    "github.com/goodsign/icu"
)

const (
    TextCatCfg = "textcat-defaultcfg/conf.txt"
    StemAlgorithmEng = "english"
    StemAlgorithmGer = "german"
    StemAlgorithmRus = "russian"
)

var (
    InputWords= []string {"supply", "name", "mein", "гора"}
    WordRx *regexp.Regexp = regexp.MustCompile(`[\d\p{L}]+`)
)

var (
    cat *libtextcat.TextCat
    detector *icu.CharsetDetector
    converter *icu.CharsetConverter
)

func search(filepath string, words []string) {
    fmt.Printf("### FILE: %s\n", filepath)

    // Read file
    bytes, err := ioutil.ReadFile(filepath)

    if err != nil {
        fmt.Println(err)
        return
    }

    // Convert to UTF-8
    converted, maxenc, err := convertToUtf8(bytes)
    
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("    Detected encoding:\n")
    fmt.Printf("        %s\n", maxenc)


    // Extract words
    docWords := WordRx.FindAllString(string(converted), -1)

    // Get possible languages
    langs, err := getPossibleLanguages(converted)
    
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("    Detected languages:\n")

    for _, l := range langs {
        fmt.Printf("        %s\n", l)
    }

    fmt.Println()
    fmt.Printf("    Found words:\n")

    // Compare stems for each possible language
    for _, dw := range docWords {
        for _, w := range words {
            for _, l := range langs {
                dwstem, err := getStem(dw, l)

                if err != nil {
                    fmt.Println(err)
                    return
                }

                wstem, err := getStem(w, l)

                if err != nil {
                    fmt.Println(err)
                    return
                }
                

                if dwstem == wstem {
                    fmt.Printf("    %s.     (Original word: '%s' Stem: '%s')\n", dw, w, dwstem)
                    break
                }
            }
        }
    }

    fmt.Println()
}

func convertToUtf8(text []byte) (converted []byte, detenctedEnc string, e error) {
    // Guess encoding
    encMatches, err := detector.GuessCharset(text)

    if err != nil {
        return nil, "", err
    }

    // Get charset with max confidence (goes first)
    maxenc := encMatches[0].Charset

    // Convert to utf-8
    converted, err = converter.ConvertToUtf8(text, maxenc)

    if err != nil {
        return nil, "", err
    }

    return converted, maxenc, nil
}

func getPossibleLanguages(text []byte) ([]string, error) {
    matches, err := cat.Classify(string(text))

    if err != nil {
        return nil, err
    }

    return matches, nil
}

func getStem(word string, language string) (string, error) {
    // Create stemmer
    stemmer, err := snowball.NewWordStemmer(language, snowball.DefaultEncoding)

    if err != nil {
        return "", err
    }
    defer stemmer.Close() 

    wordStem, err := stemmer.Stem([]byte(word))
    if err != nil {
        return "", err
    }

    return string(wordStem), nil
}

func main() {

// --------------------
// Create toolchain
// --------------------

    var err error 
    // Create textcat (Language detector)
    cat, err = libtextcat.NewTextCat(TextCatCfg)

    if err != nil {
        fmt.Println(err)
        return
    }
    defer cat.Close()

    // Create icu charset detector
    detector, err = icu.NewCharsetDetector()

    if err != nil {
        fmt.Println(err)
        return
    }
    defer detector.Close()


    // Create icu charset converter
    converter = icu.NewCharsetConverter(icu.DefaultMaxTextSize)

// --------------------
// Search words
// --------------------
    search("input-documents/eng_88591.txt", InputWords)
    search("input-documents/german_utf16le.txt", InputWords)
    search("input-documents/rus_koi8r.txt", InputWords)
}