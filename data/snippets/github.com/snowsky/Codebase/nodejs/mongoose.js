var mongoose = require('mongoose');

mongoose.connect('mongodb://localhost/superhr');

var blogSchema = {
    title:  String,
    author: String,
    body:   String,
    comments: [{ body: String, date: Date }],
    date: { type: Date, default: Date.now },
    hidden: Boolean,
    meta: {
        votes: Number,
        favs:  Number
    }
};

var Blog = mongoose.model('Blog', blogSchema);

var blog = new Blog({
    title: 'this is my blog title',
    author: 'me',
    body: 'the body of my blog. can you see that?'        
});

blog.save();

mongoose.disconnect();


//Blog.create({ 
//    title: 'another blog title', 
//    author: 'still me', 
//    body: 'the blog body again!' 
//}, function (err, small) {
//  if (err) return handleError(err);
//  // saved!
//});

