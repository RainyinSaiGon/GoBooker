import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import Providers from "./providers";
import "./globals.css";

const geistSans = Geist({ variable: "--font-geist-sans", subsets: ["latin"] });
const geistMono = Geist_Mono({ variable: "--font-geist-mono", subsets: ["latin"] });

/**
 * Root metadata — docs recommend the `template` pattern so every child page
 * automatically gets branded titles without repeating the app name.
 * https://nextjs.org/docs/app/getting-started/metadata-and-og-images
 */
export const metadata: Metadata = {
  // title.template applies to all child pages that export their own title.
  // title.default is the fallback when a child exports no title.
  title: {
    default: "GoBooker – Book smarter",
    template: "%s | GoBooker",
  },
  description:
    "GoBooker is the modern booking platform for teams and individuals. Manage users, bookings, and schedules — all in one place.",
  // metadataBase is required for resolving relative OG image paths.
  metadataBase: new URL(
    process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000"
  ),
  openGraph: {
    title: "GoBooker – Book smarter",
    description: "The modern booking platform for teams and individuals.",
    type: "website",
    locale: "en_US",
    siteName: "GoBooker",
  },
  twitter: {
    card: "summary_large_image",
    title: "GoBooker – Book smarter",
    description: "The modern booking platform for teams and individuals.",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: { index: true, follow: true },
  },
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html
      lang="en"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="h-full flex flex-col">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
