const path = require('path');
const process = require('process');
const CopyWebpackPlugin = require('copy-webpack-plugin');

/* eslint-disable import/no-extraneous-dependencies */
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});

const pathBuilder = (subpath) => path.join(process.cwd(), subpath);

module.exports = withBundleAnalyzer({
  poweredByHeader: false,
  trailingSlash: true,
  basePath: '',
  reactStrictMode: true,
  webpack: (config, { webpack, isServer }) => {
    if (!isServer) {
      config.plugins.push(
        new CopyWebpackPlugin({
          patterns: [
            {
              from: pathBuilder('node_modules/cesium/Build/Cesium/Workers'),
              to: path.join(process.cwd(), 'public/cesium/Workers'),
              info: { minimized: true },
            },
            {
              from: pathBuilder('node_modules/cesium/Build/Cesium/ThirdParty'),
              to: path.join(process.cwd(), 'public/cesium/ThirdParty'),
              info: { minimized: true },
            },
            {
              from: pathBuilder('node_modules/cesium/Build/Cesium/Assets'),
              to: path.join(process.cwd(), 'public/cesium/Assets'),
              info: { minimized: true },
            },
            {
              from: pathBuilder('node_modules/cesium/Build/Cesium/Widgets'),
              to: path.join(process.cwd(), 'public/cesium/Widgets'),
              info: { minimized: true },
            },
          ],
        }),
        new webpack.DefinePlugin({
          CESIUM_BASE_URL: JSON.stringify('/cesium'),
        })
      );
    }

    return config;
  },
  output: 'standalone',
});
