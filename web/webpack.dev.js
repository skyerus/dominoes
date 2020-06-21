const merge = require('webpack-merge')
const common = require('./webpack.common.js')
const HtmlWebPackPlugin = require('html-webpack-plugin')

module.exports = merge(common, {
  mode: 'development',
  devtool: 'inline-source-map',
  devServer: {
    historyApiFallback: true,
    noInfo: true,
    disableHostCheck: true,
    host: '0.0.0.0',
    port: 8081,
  },
  module: {
    rules: [
    ]
  },
  plugins: [
    new HtmlWebPackPlugin({
      template: "./src/index.html",
      filename: "./index.html"
    }),
  ]
})

