/*
var mongoose = require('mongoose');

var uri = 'mongodb://localhost/superhr';

mongoose.connect(uri);

var mailSchema = {
    id: String,
    threadId: String
};

var hrMail = mongoose.model('hrMail', mailSchema);

var mail = new hrMail({
    id: 'test',
    threadId: 'test' 
});

mail.save();

mail.save(function(err) {
    console.log("Cannot insert to db! " + err);
});

mongoose.disconnect();
*/

var mongoose = require('mongoose');
mongoose.connect('mongodb://localhost/superhr');
var db = mongoose.connection;
var Schema = mongoose.Schema;
db.on('error', console.error.bind(console, 'connection error:'));
db.once('open', function callback () {
    var mailSchema = new Schema({
        id: String,
        threadId: String
    }, {
        strict: false
    });
    var hrMail = mongoose.model('hrMail', mailSchema);
    var mail = new hrMail({
        id: 'test',
        threadId: 'test' 
    });
    mail.save(function(err) {
        console.log("Cannot insert to db! " + err);
    });
});
db.on('disconnected', function () {
    mongoose.disconnect();
});
