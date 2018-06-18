package main

import (
  "strings"
  "regexp"
  "sort"
  "io/ioutil"
  "github.com/reiver/go-porterstemmer"
)


func Tag_Filter(info_raw *string) (tags_filtered []string) {
  var stem string
  csv_content, err := ioutil.ReadFile("stopwords.csv")
  if (err != nil) { panic(err) }

  stopwords := strings.Split(string(csv_content), ",")
  tags_raw := strings.Fields(*info_raw)
  tags_filtered_map := make(map[string]int)
  sort.Sort(sort.StringSlice(tags_raw))

  i_stopwords := 0
  i_tags_raw := 0
  i_tags_filtered_map := 1 // shift by 1 to stop comparing index by 0 (*A)

  for (i_stopwords < len(stopwords)) && (i_tags_raw < len(tags_raw)) {
    if (tags_raw[i_tags_raw] == stopwords[i_stopwords]) {
      i_tags_raw++
    } else if (tags_raw[i_tags_raw] > stopwords[i_stopwords]) {
      i_stopwords++
    } else {
      stem = Tag_Stemmed(tags_raw[i_tags_raw])
      if (tags_filtered_map[stem] == 0) && (tags_raw[i_tags_raw] != "") {
        tags_filtered_map[stem] = i_tags_filtered_map
        i_tags_filtered_map++
      }
      i_tags_raw++
    }
  }
  for (i_tags_raw < len(tags_raw)) { // get the remaining tags
    stem = Tag_Stemmed(tags_raw[i_tags_raw])
    if (tags_filtered_map[stem] == 0) && (tags_raw[i_tags_raw] != "") {
      tags_filtered_map[stem] = i_tags_filtered_map
      i_tags_filtered_map++
    }
    i_tags_raw++
  }

  tags_filtered = make([]string, len(tags_filtered_map), len(tags_filtered_map))
  for tag, index := range tags_filtered_map {
    tags_filtered[index-1] = tag // remove the shift of 1 (*A)
  }
  return tags_filtered
}


func Tag_CommonalityScore (tags_a []string, tags_b []string) (score int) {
  score = 0
  i_tags_a := 0
  i_tags_b := 0
  for (i_tags_a < len(tags_a)) && (i_tags_b < len(tags_b)) {
    if (tags_b[i_tags_b] == tags_a[i_tags_a]) {
      score++
      i_tags_a++
      i_tags_b++
    } else if (tags_b[i_tags_b] > tags_a[i_tags_a]) {
      i_tags_a++
    } else {
      i_tags_b++
    }
  }
  return score
}


func Tag_Stemmed (tag_raw string) (tag_stemmed string) {
  regexp_lowercase, err := regexp.Compile("[^a-z]")
  if err != nil { panic(err) }
  tag_stemmed = regexp_lowercase.ReplaceAllString(strings.ToLower( tag_raw ), "")
  tag_stemmed = porterstemmer.StemString(tag_stemmed)
  return tag_stemmed
}


