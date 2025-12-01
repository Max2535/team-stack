import "./../styles/globals.css";
import type { ReactNode } from "react";

export const metadata = {
  title: "Team Stack",
  description: "Next.js + Go Fiber template"
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
