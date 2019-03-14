const webpack = require('webpack');
const path = require('path');
const VueLoaderPlugin = require('vue-loader/lib/plugin');
const UglifyJsPlugin = require('uglifyjs-webpack-plugin');
const BrowserSyncPlugin = require('browser-sync-webpack-plugin')


module.exports = {
    devtool: false,
    entry: {
        main: './resources/js/main.js',
        admin: './resources/js/admin.js',
    },
    output: {
        path: path.resolve(__dirname, 'static/js'),
        filename: '[name].js'
    },
    resolve: {
        alias: {
            'vue$': 'vue/dist/vue.esm.js'
        },
        extensions: ['*', '.js', '.vue', '.json']
    },
    module: {
        rules: [
            {
                test: /\.vue$/,
                loader: 'vue-loader',
            },
            {
                test: /\.js$/,
                loader: 'babel-loader'
            },
            {
                test: /\.css$/,
                use: [
                    'vue-style-loader',
                    'css-loader',
                ]
            },
            // {
            //     test: /\.scss$/,
            //     use: [
            //         'vue-style-loader',
            //         'css-loader',
            //         'sass-loader'
            //     ]
            // },
        ],
    },
    optimization: {
        // minimizer: [
        //     new UglifyJsPlugin({
        //         cache: true,
        //         parallel: true,
        //         sourceMap: true,
        //         uglifyOptions: {
        //             compress: {
        //                 drop_console: true
        //             }
        //         }
        //     }),
        // ]
    },
    plugins: [
        // make sure to include the plugin!
        new VueLoaderPlugin(),
        new webpack.SourceMapDevToolPlugin({}),
        new BrowserSyncPlugin({
            // browse to http://localhost:3000/ during development,
            // ./public directory is being served
            host: 'localhost',
            port: 3000,
            server: { baseDir: ['./'] }
        })
    ],
};