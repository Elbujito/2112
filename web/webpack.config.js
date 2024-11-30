const path = require('path');
const webpack = require('webpack');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');
const NodePolyfillPlugin = require('node-polyfill-webpack-plugin');

const cesiumSource = path.resolve(__dirname, 'node_modules/cesium/Source');
const cesiumWorkers = path.resolve(cesiumSource, '../Build/Cesium/Workers');

module.exports = {
    context: __dirname,
    mode: 'development',
    entry: {
        app: './src/index.js', // Entry point for the application
    },
    output: {
        filename: '[name].js',
        path: path.resolve(__dirname, 'dist'),
        sourcePrefix: '', // Cesium requires an empty sourcePrefix
    },
    amd: {
        toUrlUndefined: true, // Required by Cesium
    },
    module: {
        rules: [
            {
                test: /\.css$/,
                use: ['style-loader', 'css-loader'], // Handle CSS files
            },
            {
                test: /\.(png|gif|jpg|jpeg|svg|xml|json)$/,
                use: ['url-loader'], // Handle assets
            },
        ],
    },
    plugins: [
        new HtmlWebpackPlugin({
            template: 'src/index.html', // Use your custom HTML template
        }),
        new NodePolyfillPlugin(), // Polyfills for Node.js modules
        new CopyWebpackPlugin({
            patterns: [
                { from: cesiumWorkers, to: 'Workers' },
                { from: path.join(cesiumSource, 'Assets'), to: 'Assets' },
                { from: path.join(cesiumSource, 'Widgets'), to: 'Widgets' },
            ],
        }),
        new webpack.DefinePlugin({
            CESIUM_BASE_URL: JSON.stringify(''), // Define base URL for Cesium
        }),
    ],
    devServer: {
        static: {
            directory: path.resolve(__dirname, 'dist'), // Serve files from 'dist'
        },
        compress: true,
        port: 8082, // Port for the dev server
        hot: true, // Enable hot module replacement
        open: true, // Automatically open the browser
    },
    resolve: {
        alias: {
            cesium: cesiumSource, // Resolve Cesium source directory
        },
        extensions: ['.js', '.json'], // Resolve these extensions
    },
};
