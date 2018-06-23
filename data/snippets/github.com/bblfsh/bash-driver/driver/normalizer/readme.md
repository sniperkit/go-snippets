# Problems

## Tokens are the original full text, not just the symbol

For code like `$A` the bash parser returns a leaf like this:

```
{
  "children": [],
  "elementType": "[Bash] variable",
  "startOffset": 0,
  "text": "$A",
  "textLength": 2
}
```

The text is "$A", not "A" as you would expect from parsers of other
languages like Java.

## Tokens of parents hold the text of their children

For code like `$A` the bash parser returns a tree like this:

```
{
  "children": [
  {
    "children": [
    {
      "children": [
      {
        "children": [
        {
          "children": [],
          "elementType": "[Bash] variable",
          "startOffset": 0,
          "text": "$A",
          "textLength": 2
        }
        ],
        "elementType": "var-use-element",
        "startOffset": 0,
        "text": "$A",
        "textLength": 2
      }
      ],
      "elementType": "[Bash] combined word",
      "startOffset": 0,
      "text": "$A",
      "textLength": 2
    }
    ],
    "elementType": "[Bash] generic bash command",
    "startOffset": 0,
    "text": "$A",
    "textLength": 2
  }
  ],
  "elementType": "simple-command",
  "startOffset": 0,
  "text": "$A",
  "textLength": 2
},
```

Note how every node contains the text of their children.

This is particularly painfull in the case of the "FILE" node that contains the
full sourced text of the file being analysed.

## We need a predicate to match then nth child

In bash, the IfCondition, IfBody and IfElse are identified as the 3, 8, and 11
children of the a node idenfied as "if shellcommand", for instance:

```bash
if /bin/false; then /bin/true; fi
```

Has the following native:

```
"if shellcommand"
|__ "[Bash] if"
|__ "WHITE_SPACE"
|__ "simple-command  <-- this is the if condition, but you cannot
     |                   know that from its element type, you need its position
     |                   in the array.
     |__ ...
     |__ ...
     |__ ...
|__ "[Bash] ;"
|__ "WHITE_SPACE"
|__ "[Bash] then"
|__ "WHITE_SPACE"
|__ "logical block"  <-- this is the if body, but you cannot know that from
     |                   its element type, you need its position in the array.
     |__ ...
     |__ ...
     |__ ...
|__ "[Bash] ;"
|__ "WHITE_SPACE"
|__ "[Bash] fi"

Ideally I would like to say, On("if shellcommand"), the 3rd child is the
condition and the 8th thing is the body.
