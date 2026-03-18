import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // API proxy to Go backend (optional: use in dev to avoid CORS)
  async rewrites() {
    return [
      { source: "/api/:path*", destination: "http://localhost:8080/:path*" },
    ];
  },
};

export default nextConfig;
