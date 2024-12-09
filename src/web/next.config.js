/** @type {import('next').NextConfig} */
const path = require('path');
const process = require('process');
const CopyWebpackPlugin = require('copy-webpack-plugin');

// Bundle analyzer for debugging build size
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true',
});

const pathBuilder = (subpath) => path.join(process.cwd(), subpath);

const nextConfig = withBundleAnalyzer({
  poweredByHeader: false,
  trailingSlash: true,
  reactStrictMode: true,
  swcMinify: true,
  basePath: process.env.NEXT_PUBLIC_BASE_PATH,
  assetPrefix: process.env.NEXT_PUBLIC_BASE_PATH,
  images: {
    domains: [
      'images.unsplash.com',
      'i.ibb.co',
      'scontent.fotp8-1.fna.fbcdn.net',
    ],
    unoptimized: true,
  },
  webpack: (config, { webpack, isServer }) => {
    if (!isServer) {
      // Add Cesium-specific configurations
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
          CESIUM_BASE_URL: JSON.stringify('/cesium'), // Set Cesium base URL
        })
      );
    }

    return config;
  },
  output: 'standalone',
});

module.exports = nextConfig;
