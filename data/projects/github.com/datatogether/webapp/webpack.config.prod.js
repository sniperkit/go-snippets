'use strict';

var path = require('path');
var webpack = require('webpack');

var ENV = {
  __BUILD__ : {
    PRODUCTION : JSON.stringify(true),
    DEVELOP : JSON.stringify(false),
    STAGING : JSON.stringify(false),

    BASE_URL : JSON.stringify("https://archivers.co"),
    API_URL : JSON.stringify("https://archivers.co"),
    USERS_API_URL : JSON.stringify("https://ident.archivers.space"),
    COVERAGE_API_URL : JSON.stringify("https://coverage.archivers.space"),
    WEBSOCKET_URL : JSON.stringify("wss://archivers.co/ws"),
    STATIC_ASSETS_URL : JSON.stringify("http://archivers.co"),
    // SEGMENT_KEY : JSON.stringify("FwUGnmKzryJpDpApVzdQy9rmwwWSiK1M"),
    SEGMENT_KEY : JSON.stringify(""),
    SEGMENT_APP_PROPERTIES: JSON.stringify('{ "app_version": "2.0.0beta1" }'),
  }
};

function exitOneFail() {
  this.plugin("done", function(stats) {
    if (stats.compilation.errors && stats.compilation.errors.length) {
      console.log(stats.compilation.errors);
      process.exit(1);
    }
  });
}

module.exports = {
  devtool : 'cheap-module-source-map',
  entry: {
    'app' : './src/js/index.js',
    'vendor' : [ 'react', 'react-dom' ]
  },
  output: {
    path: path.join(__dirname, 'dist'),
    publicPath : "https://s3.amazonaws.com/static.qri.io/js/archivers/co/",
    filename: '[name].min.js',
    chunkFilename: "[name].chunk.min.js"
  },
  plugins: [
    // new webpack.optimize.OccurenceOrderPlugin(),
    new webpack.DefinePlugin(ENV),
    new webpack.DefinePlugin({ 'process.env': { 'NODE_ENV': JSON.stringify('production') } }),
    exitOneFail,
    new webpack.optimize.UglifyJsPlugin({minimize: true}),
  ],
  module: {
    rules: [
      {
        test: /\.js?$/,
        loaders: ['babel-loader'],
        exclude: /node_modules/
      },
      {
        test: /\.scss$/,
        loaders: [
          "style-loader",
          "css-loader",
          "sass-loader"
        ]
      },
      {
        test: /\.css$/,
        use: [
          "css-loader",
        ]
      }
    ]
  }
};
