import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "SecScan",
  description: "Web security scanning dashboard for the Seri 3 final project"
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="tr">
      <body>{children}</body>
    </html>
  );
}

