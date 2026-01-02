import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { AuthProvider } from "@/lib/auth";
import Header from "@/components/Header";
import Footer from "@/components/Footer";
import StarBackground from "@/components/StarBackground";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Matthew's Galaxy | Matei-Alexandru Dinu",
  description: "Exploring the frontiers of technology. A blog by Matei-Alexandru Dinu, Software Engineer with experience at Google, Snowflake, and NXP Semiconductors.",
  keywords: ["software engineering", "cloud computing", "golang", "kubernetes", "tech blog"],
  authors: [{ name: "Matei-Alexandru Dinu" }],
  openGraph: {
    title: "Matthew's Galaxy",
    description: "Exploring the frontiers of technology",
    url: "https://matthewsgalaxy.com",
    siteName: "Matthew's Galaxy",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <AuthProvider>
          <StarBackground />
          <Header />
          <main style={{ paddingTop: 'var(--header-height)', minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
            {children}
          </main>
          <Footer />
        </AuthProvider>
      </body>
    </html>
  );
}
