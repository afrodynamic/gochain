import { Health } from '@/app/(features)/(dashboard)/Health';
import { KeyGenerator } from '@/app/(features)/(dashboard)/KeyGenerator';
import { Wallet } from '@/app/(features)/(dashboard)/Wallet';
import { Navbar } from '@/app/_components/layout/Navbar';
import Masonry from '@mui/lab/Masonry';
import { Box } from '@mui/material';

export default function HomePage() {
  return (
    <Box className="min-h-svh">
      <Navbar />

      <main className="w-full">
        <Masonry
          columns={{ xs: 1, md: 2, lg: 3 }}
          spacing={2}
          className="w-full p-2 max-w-7xl mx-auto"
        >
          <Box className="max-lg:flex max-lg:justify-center">
            <Health />
          </Box>

          <Box className="max-lg:flex max-lg:justify-center">
            <Wallet />
          </Box>

          <Box className="max-lg:flex max-lg:justify-center">
            <KeyGenerator />
          </Box>
        </Masonry>
      </main>
    </Box>
  );
}
