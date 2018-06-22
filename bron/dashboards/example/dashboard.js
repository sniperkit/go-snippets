var blessed = require('blessed')
  , contrib = require('../../index')

var screen = blessed.screen()

//create layout and widgets
var grid = new contrib.grid({rows: 1, cols: 2})
var grid1 = new contrib.grid({rows: 1, cols: 1})
grid1.set(0, 0, contrib.line,
  { style:
    { line: "blue"
    , text: "white"
    , baseline: "black"}
  , label: 'Number of Files'
  , maxY: 100})
var grid2 = new contrib.grid({rows: 2, cols: 1})
grid2.set(0, 0, contrib.bar,
  { label: 'Lines per Language'
  , barWidth: 3
  , barSpacing: 6
  , xOffset: 0
  , maxHeight: 5})
grid2.set(1, 0, contrib.table,
  { keys: true
  , fg: 'green'
  , label: 'Number of Commits per Author'
  , columnSpacing: [24, 10]})
var grid3 = new contrib.grid({rows: 3, cols: 1})
grid3.set(0, 0, contrib.line,
  { style:
    { line: "red"
    , text: "white"
    , baseline: "black"}
  , label: 'Number of Lines'
  , maxY: 1000})
grid3.set(1, 0, grid2)
grid3.set(2, 0, grid1)
var grid4 = new contrib.grid({rows: 2, cols: 1})
grid4.set(0, 0, contrib.line,
  { showNthLabel: 5
  , maxY: 100
  , label: 'Number of Languages'})
grid4.set(1, 0, contrib.line,
  { style:
    { line: "green"
    , text: "white"
    , baseline: "black"}
  , label: 'Number of Authors'
  , maxY: 10})
grid.set(0, 0, grid4)
grid.set(0, 1, grid3)
grid.applyLayout(screen)

var numLanguages = grid4.get(0, 0)
var numLines = grid3.get(0, 0)
var numAuthors = grid4.get(1, 0)
var numFiles = grid1.get(0, 0)
var commitsAuthor = grid2.get(1,0)
var linesLanguage = grid2.get(0, 0)
var languages = $languages
var languageLines = $languageLines
var authors = $authors
var numLanguagesData = $numLanguagesData
var numLinesData = $numLinesData
var numAuthorsData = $numAuthorsData
var numFilesData = $numFilesData

//set data on bar chart
function fillBar() {
  linesLanguage.setData({titles: languages, data: languageLines})
}
fillBar()

//set data for table
function generateTable() {
   commitsAuthor.setData({headers: ['Author', 'Commits'], data: authors})
}
generateTable()
commitsAuthor.focus()

//set line charts data
function setLineData(data, line) {
  line.setData(data.x, data.y)
}
setLineData(numLanguagesData, numLanguages)
setLineData(numLinesData, numLines)
setLineData(numAuthorsData, numAuthors)
setLineData(numFilesData, numFiles)

screen.key(['escape', 'q', 'C-c'], function(ch, key) {
  return process.exit(0);
});

screen.render()
