var path = require('path');
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var StyleLintWebpackPlugin = require('stylelint-webpack-plugin');
var CleanWebpackPlugin = require('clean-webpack-plugin');
var ExtractTextWebpackPlugin = require('extract-text-webpack-plugin');

module.exports = {
  entry: [
    './src/index.tsx',
  ],

  output: {
    path: path.resolve(__dirname, 'dist.prod'),
    filename: 'bundle.js',
  },

  plugins: [
    new CleanWebpackPlugin([
      'dist.prod',
      'src/**/*.d.ts',
    ]),
    new webpack.DefinePlugin({
      '__DEV__': false,
      'process.env.NODE_ENV': JSON.stringify('production'),
    }),
    new StyleLintWebpackPlugin({
      files: [
        '**/*.less',
      ],
      syntax: 'less',
      config: {
        'extends': 'stylelint-config-standard',
      },
    }),
    new ExtractTextWebpackPlugin({
      filename: 'style.css',
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.prod.html',
      hash: true,
    }),
    new webpack.optimize.UglifyJsPlugin({
      sourceMap: true,
    }),
  ],

  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: ['.ts', '.tsx', '.js'],
  },

  module: {
    rules: [
      {
        test: /\.less$/,
        include: path.resolve(__dirname, 'src'),
        use: ExtractTextWebpackPlugin.extract({
          fallback: 'style-loader',
          use: [
            {
              loader: 'css-loader',
              options: {
                modules: true,
                sourceMap: true,
                camelCase: true,
                namedExport: true,
                localIdentName: '__[name]--[local]',
              },
            },
            'autoprefixer-loader',
            'less-loader',
          ],
        }),
      },
      {
        test: /\.tsx?$/,
        include: path.resolve(__dirname, 'src'),
        use: [
          'ts-loader',
        ],
      },
      {
        test: /\.tsx?$/,
        include: path.resolve(__dirname, 'src'),
        enforce: 'pre',
        use: [
          {
            loader: 'tslint-loader',
            options: {
              emitErrors: true,
            },
          },
        ],
      },
      {
        test: /\.svg$/,
        include: path.resolve(__dirname, 'src'),
        use: [
          'raw-loader',
        ],
      },
    ],
  },

  // When importing a module whose path matches one of the following, just
  // assume a corresponding global variable exists and use that instead.
  // This is important because it allows us to avoid bundling all of our
  // dependencies, which allows browsers to cache those libraries between builds.
  externals: {
    'es5-shim': 'es5',
    'es6-shim': 'es6',
    'react': 'React',
    'react-dom': 'ReactDOM',
    'react-router': 'ReactRouter',
    'mobx': 'mobx',
    'mobx-react': 'mobxReact',
  },
};
