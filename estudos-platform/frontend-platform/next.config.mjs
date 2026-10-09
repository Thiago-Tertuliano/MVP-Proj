/** @type {import('next').NextConfig} */
const nextConfig = {
  // Sem domínio próprio, front (Vercel) e API (Render) ficam em sites diferentes e o navegador
  // descartaria os cookies de sessão. Com NEXT_PUBLIC_API_URL="/", o browser chama /api/v1 no
  // próprio front e este proxy repassa para API_URL — os cookies continuam first-party.
  async rewrites() {
    const api = process.env.API_URL?.replace(/\/+$/, "");
    return api ? [{ source: "/api/v1/:path*", destination: `${api}/api/v1/:path*` }] : [];
  },
};

export default nextConfig;
