// Fetch latest 10 emails and show the snippet
 
var Gmail = require('node-gmail-api')
  , gmail = new Gmail("WXyGBVBPDiDzMdspWP41oW8y")
  , s = gmail.messages('label:inbox', {max: 10})
 
s.on('data', function (d) {
  console.log(d.snippet)
})
