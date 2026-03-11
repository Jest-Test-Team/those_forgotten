import type { Metadata } from "next";
import { Noto_Sans_TC, Noto_Serif_TC } from "next/font/google";
import { PwaRegister } from "@/components/pwa-register";
import "./globals.css";

const notoSans = Noto_Sans_TC({
  variable: "--font-body",
  subsets: ["latin"],
});

const notoSerif = Noto_Serif_TC({
  variable: "--font-title",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  metadataBase: new URL("https://customs-auction.local"),
  title: {
    default: "海關標售雷達",
    template: "%s | 海關標售雷達",
  },
  description: "追蹤全台海關標售公告、教學內容、歷史成交與代標媒合的 zh-TW Web PWA。",
  manifest: "/manifest.webmanifest",
  openGraph: {
    title: "海關標售雷達",
    description: "官方連結、風險警示、推播訂閱與實戰教學整合在同一個平台。",
    type: "website",
    locale: "zh_TW",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="zh-Hant">
      <body
        className={`${notoSans.variable} ${notoSerif.variable} antialiased`}
      >
        <PwaRegister />
        {children}
      </body>
    </html>
  );
}
