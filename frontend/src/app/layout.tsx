import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { AppRouterCacheProvider } from "@mui/material-nextjs/v16-appRouter";
import "./globals.css";
import { MSWProvider } from "@/components/MSWProvider";
import { MuiThemeProvider } from "@/components/MuiThemeProvider";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Trabelle",
  description: "Plan and share your travel itinerary",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}>
      <body className="h-full flex flex-col overflow-hidden">
        <AppRouterCacheProvider>
          <MuiThemeProvider>
            {process.env.NODE_ENV === "development" ? (
              <MSWProvider>{children}</MSWProvider>
            ) : (
              children
            )}
          </MuiThemeProvider>
        </AppRouterCacheProvider>
      </body>
    </html>
  );
}
