import type { Metadata } from "next";
import localFont from "next/font/local";

import { SessionHeader } from "@/components/relp/SessionHeader";
import { RewardHost } from "@/components/roadmap/RewardHost";
import { Toaster } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { SessionProvider } from "@/lib/session";
import "./globals.css";

const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
  weight: "100 900",
});
const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
  weight: "100 900",
});

export const metadata: Metadata = {
  title: {
    default: "Relp! — Plataforma de Aprendizado",
    template: "%s · Relp!",
  },
  description: "Trilhas interativas de tecnologia com progresso personalizado, anotações e busca inteligente",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR">
      <body className={`${geistSans.variable} ${geistMono.variable} min-h-screen font-sans`}>
        <SessionProvider>
          <TooltipProvider delayDuration={200}>
            <a
              href="#conteudo"
              className="sr-only focus:not-sr-only focus:fixed focus:left-3 focus:top-3 focus:z-50 focus:rounded-md focus:bg-primary focus:px-4 focus:py-2 focus:text-sm focus:font-medium focus:text-primary-foreground"
            >
              Pular para o conteúdo
            </a>
            <SessionHeader />
            <main id="conteudo" tabIndex={-1} className="mx-auto max-w-5xl px-4 py-8 focus:outline-none">
              {children}
            </main>
            <RewardHost />
            <Toaster />
          </TooltipProvider>
        </SessionProvider>
      </body>
    </html>
  );
}
