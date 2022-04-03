const { defineConfig } = require('@vue/cli-service')
const fs = require('fs');

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    https: true,
    server: {
      type: 'https',
      options: {
        key: fs.readFileSync('server.key'),
        cert: fs.readFileSync('server.crt')
      },
    },
    allowedHosts: [
      'localhost.hedwig100.tk',
    ],
  },
  configureWebpack: {
    resolve: {
      fallback: {
        https: require.resolve("https-browserify"),
        http: require.resolve("stream-http"),
        url: require.resolve("url/"),
      }
    }
  }
})
