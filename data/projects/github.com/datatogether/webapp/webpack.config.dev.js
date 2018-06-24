'use strict';

var path = require('path');
var webpack = require('webpack');

var ENV = {
  __BUILD__ : {
    PRODUCTION : JSON.stringify(false),
    DEVELOP : JSON.stringify(true),

    BASE_URL : JSON.stringify("http://localhost:3000"),
    BASE_URL_PORT : JSON.stringify("3000"),
    API_URL : JSON.stringify("http://localhost:3000"),
    API_PORT : JSON.stringify("3000"),
    USERS_API_URL : JSON.stringify("http://localhost:3100"),
    USERS_API_PORT : JSON.stringify("3100"),
    COVERAGE_API_URL : JSON.stringify("http://localhost:3200"),
    COVERAGE_API_PORT : JSON.stringify("3200"),
    WEBSOCKET_URL : JSON.stringify("ws://localhost:3000/ws"),
    WEBSOCKET_PORT : JSON.stringify("3000"),
    STATIC_ASSETS_URL : JSON.stringify("http://localhost:3000"),
    SEGMENT_KEY : JSON.stringify(""),
    SEGMENT_APP_PROPERTIES: JSON.stringify('{ "app_version": "2.0.0beta1" }'),
  }
};

module.exports = {
  devtool: 'cheap-module-eval-source-map',
  entry: [
    'webpack-hot-middleware/client?path=//localhost:4000/__qri_io',
    './src/js/index'
  ],
  output: {
    // path: path.join(__dirname, 'dist'),
    path : __dirname,
    filename: 'bundle.js',
    publicPath: '//localhost:4000/static/'
  },
  // resolve: {
  //   alias: {
  //     'react-ace':  path.join(__dirname, '..', 'src', 'ace.jsx')
  //   }
  // },
  plugins: [
    new webpack.DefinePlugin(ENV),
    // new webpack.optimize.OccurenceOrderPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin()
  ],
  module: {
    rules: [
      {
        test: /\.js?$/,
        use: ['babel-loader'],
        exclude: /node_modules/,
      },
      {
        test: /\.scss$/,
        use: [
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