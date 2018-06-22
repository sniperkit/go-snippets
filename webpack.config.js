module.exports = {
    entry:{
        main : "./web/js/entry.js"
    },
    output: {
        path: "web/lib/",
        filename: "bundle.js"
    },
    module: {
      loaders: [
        { test: /\.js$/, exclude: /node_modules/, loader: "babel-loader"}
      ]
    }
};