const HtmlWebpackPlugin = require('html-webpack-plugin');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = {
    mode: 'development',
    plugins: [
        new HtmlWebpackPlugin({
            template: 'public/index.html'
        }),
        new CopyPlugin({
            patterns: [
                { from: "../sim/out/boids.wasm", to: "game/boids.wasm" },
                { from: "./public", to: ".", globOptions: { ignore: '**/public/index.html' } },
            ],
        }),
    ]
}