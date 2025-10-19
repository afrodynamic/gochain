import { firaCode, inter } from '@/app/_components/theme';
import { ThemeRegistry } from '@/app/_components/ThemeRegistry';
import { MetamaskProvider } from '@/app/_utils/MetamaskProvider';
import { ReactQueryProvider } from '@/app/_utils/ReactQueryProvider';
import type { Metadata } from 'next';
import { ReactNode } from 'react';
import './globals.css';

export const metadata: Metadata = {
  title: 'Gochain',
  description: 'The ultimate blockchain platform.',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${inter.variable} ${firaCode.variable} antialiased`}>
        <MetamaskProvider>
          <ReactQueryProvider>
            <ThemeRegistry>{children}</ThemeRegistry>
          </ReactQueryProvider>
        </MetamaskProvider>
      </body>
    </html>
  );
}
