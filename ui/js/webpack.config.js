const path = require('path');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');

module.exports = {
  watch: true,
  watchOptions: {
    ignored: /node_modules/
  },
  mode: 'development',
  //devtool: false,
  entry: {
    default: './src/default',
    nav: './src/nav/nav',
    timeline: './src/timeline/timeline',
    index: './src/index/index',
    login: './src/auth/login/login',
    signup: './src/auth/signup/signup',
    userMenu: './src/user-menu/user-menu',
    createRoom: './src/create-room/create-room',
    welcome: './src/welcome/welcome',
    popup: './src/popup/popup',
  },
  module: {
    rules: [
      {
        test: /\.m?js$/,
        exclude: /node_modules/,
        use: {
          loader: "babel-loader",
          options: {
            presets: ['@babel/preset-env']
          }
        }
      },
      {
        test: /\.(html|svelte)$/,
        use: {
            loader: 'svelte-loader',
            options: {
              hydratable: true,
            },
          },
      }
    ]
  },
  resolve: {
    alias: {
      svelte: path.resolve('node_modules', 'svelte'),
      vue: 'vue/dist/vue.min.js',
    },
    extensions: ['.mjs', '.js', '.svelte', 'vue' ],
    mainFields: ['svelte', 'browser', 'module', 'main']
  },
  output: {
    filename: '[name].[contenthash].js',
    chunkFilename: '[name].[contenthash].js',
    path: path.resolve(__dirname, '../../static/js/'),
    publicPath: '/static/js/'
  },
  optimization: {
    minimize: false,
    usedExports: false,
    moduleIds: 'deterministic',
    removeAvailableModules: true,
    flagIncludedChunks: true,
    minimizer: [
      new TerserPlugin({
        parallel: true,
      }),
    ],
  },
  plugins: [
      new CleanWebpackPlugin(),
  ],
};

