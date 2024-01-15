/** @type {import('next').NextConfig} */
const nextConfig = {
    output: 'export',  // generate an 'out' folder for hosting as a static site
    trailingSlash: true, // Fixes issue where refreshing page caused a 404
}

module.exports = nextConfig
