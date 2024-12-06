/** @type {import('next').NextConfig} */
const nextConfig = {
    output: "standalone",
    reactStrictMode: false,
    images: {
        remotePatterns: [{
            protocol: 'https',
            hostname: 'assets.aceternity.com',
            port: '',
            pathname: '/**',
        }]
    },
    crossOrigin: 'use-credentials',
    devIndicators: {
        buildActivityPosition: 'bottom-right',
    },
};

export default nextConfig;
