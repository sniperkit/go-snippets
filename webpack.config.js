const webpack = require('webpack'),
    path = require('path'),
    HtmlWebpackPlugin = require('html-webpack-plugin'),
    ExtractTextPlugin = require('extract-text-webpack-plugin'),
    CleanPlugin = require('clean-webpack-plugin')

const isDev = process.env.NODE_ENV === 'development'


const config = {
    entry: {
        tab: path.join(__dirname, 'app/tab.js'),
        popup: path.join(__dirname, 'app/popup.js')
    },
    output: {
        filename: '[name].bundle.js',
        path: path.join(__dirname, './build')
    },
    module: {
        rules: [{
            test: /\.vue$/,
            loader: 'vue-loader'
        }, {
            test: /\.js$/,
            use: {
                loader: 'babel-loader'
            },
            exclude: /node_modules/
        }, {
            test: /\.(scss|sass)$/,
            use: ExtractTextPlugin.extract({
                fallback: "style-loader",
                publicPath: '../../',
                use: [{
                        loader: 'css-loader',
                        options: {
                            importLoaders: 2,
                            minimize: true
                        }
                    },
                    'sass-loader'
                ]
            })
        }, {
            test: /\.(png|jpg|gif)$/,
            use: [{
                loader: 'url-loader',
                options: {
                    name: 'img/[name].[hash:4].[ext]',
                    limit: 1024
                }
            }]
        }, {
            test: /\.(woff|svg|eot|ttf)\??.*$/,
            use: {
                loader: 'file-loader',
                options: {
                    name: 'fonts/[name].[ext]'
                }
            }
        }]
    },
    plugins: [
        new HtmlWebpackPlugin({
            title: 'New Tab',
            chunks: ['tab'],
            filename: 'tab.html'
        }),
        new HtmlWebpackPlugin({
            title: 'Popup',
            chunks: ['popup'],
            filename: 'popup.html'
        }),
        new ExtractTextPlugin('style.css')
    ]
}

if (isDev) {
    config.devtool = '#cheap-module-eval-source-map'

} else {
    config.plugins.push(
        new CleanPlugin('./build', {
            // 一般不更新图标 所以不删除fonts
            //
            exclude: ['manifest.json','static', 'fonts']
        }), 
        new webpack.optimize.OccurrenceOrderPlugin(),
        new webpack.optimize.UglifyJsPlugin())
}

module.exports = config