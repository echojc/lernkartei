var sleep = require('sleep');
var path = require('path');
var url = require('url');
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var StyleLintWebpackPlugin = require('stylelint-webpack-plugin');

module.exports = {
  entry: [
    'react-hot-loader/patch',
    './src/index.tsx',
  ],

  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist.dev'),
    publicPath: '/',
  },

  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NamedModulesPlugin(),
    new webpack.DefinePlugin({
      __DEV__: true,
    }),
    new StyleLintWebpackPlugin({
      files: ['**/*.less'],
      syntax: 'less',
      config: {
        'extends': 'stylelint-config-standard',
      },
    }),
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      hash: true,
    }),
  ],

  devtool: 'eval-source-map',

  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: ['.ts', '.tsx', '.js'],
  },

  module: {
    rules: [
      {
        test: /\.less$/,
        include: path.resolve(__dirname, 'src'),
        use: [
          'style-loader',
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
      },
      {
        test: /\.tsx?$/,
        include: path.resolve(__dirname, 'src'),
        use: [
          'react-hot-loader/webpack',
          'ts-loader',
        ],
      },
      {
        test: /\.tsx?$/,
        include: path.resolve(__dirname, 'src'),
        enforce: 'pre',
        use: [
          'tslint-loader',
        ],
      },
      {
        test: /\.(svg|json)$/,
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

  devServer: {
    hot: true,
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        pathRewrite: function(path, req) {
          const u = url.parse(path)
          let localPath = u.pathname.replace(/^\/api/, '');
          if (/^\/skill\/[0-9]+\/image$/.test(localPath)) {
            console.log(`blob ${localPath}`);
          } else {
            const method = req.method;
            if (/^\/location\/[0-9\.-]+\/[0-9\.-]+/.test(localPath)) {
              console.log(`rewrite ${localPath} -> /location`);
              localPath = '/location';
            }
            localPath += `.${method.toLowerCase()}.json`;
            sleep.msleep(300);
          }
          console.log(`mock: ${req.method} ${localPath} ${u.query ? `(query: ${u.query})` : ''}`);
          req.method = 'GET';
          return `/mock-server${localPath}`;
        },
      },
    },
  },
};
