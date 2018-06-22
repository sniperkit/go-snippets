# Similarity
REST API for search by Similarity using Go

## Endpoint example

```
POST http://localhost:8080/similarity
```

## Request Payload example
```json
{
   "dictionary":[
      "Afghanistan",
      "Albania",
      "Algeria",
      "Belarus",
      "Botswana",
      "Bouvet Island",
      "Brazil",
      "British Antarctic Territory",
      "Croatia",
      "Cuba",
      "Cyprus",
      "Czech Republic",
      "East Germany",
      "Ecuador",
      "Egypt",
      "El Salvador",
      "Equatorial Guinea",
      "French Southern Territories",
      "French Southern and Antarctic Territories",
      "Gabon",
      "Hungary",
      "Iceland",
      "India",
      "Indonesia",
      "Iran",
      "Iraq",
      "Ireland",
      "Isle of Man",
      "Israel",
      "Italy",
      "Jamaica",
      "Japan",
      "Marshall Islands",
      "Netherlands",
      "Netherlands Antilles",
      "Neutral Zone",
      "New Caledonia",
      "New Zealand",
      "Nicaragua",
      "Niger",
      "Nigeria",
      "Niue",
      "Norfolk Island",
      "North Korea",
      "South Africa",
      "South Georgia and the South Sandwich Islands",
      "Yemen",
      "Zambia",
      "Zimbabwe",
      "Åland Islands"
   ],
   "input":"zabs",
   "level":"LOW"
}
```

## Response Payload example

```json
{
   "resultCode":"SUCCESS",
   "resultData":{
      "results":[
         "Zambia",
         "Zimbabwe"
      ]
   },
   "errors":null
}
```

## Json schema for request payload

```json
{
   "type":"object",
   "definitions":{

   },
   "$schema":"http://json-schema.org/draft-07/schema#",
   "properties":{
      "dictionary":{
         "$id":"/properties/dictionary",
         "type":"array",
         "items":{
            "$id":"/properties/dictionary/items",
            "type":"string",
            "default":""
         }
      },
      "input":{
         "$id":"/properties/input",
         "type":"string",
         "default":""
      },
      "level":{
         "$id":"/properties/level",
         "type":"string",
         "default":"",
         "enum":[
            "EXTRA_LOW",
            "LOW",
            "MEDIUM",
            "HIGH",
            "EXTRA_HIGH"
         ]
      }
   }
}
```

## Docker Hub image

https://hub.docker.com/r/samuellobato/similarity/