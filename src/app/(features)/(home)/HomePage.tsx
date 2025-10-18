import { Health } from '@/app/(features)/(dashboard)/Health';
import { Navbar } from '@/app/_components/layout/Navbar';
import { Box } from '@mui/material';

export default function HomePage() {
  return (
    <Box className="flex flex-col">
      <Navbar />

      <main className="flex flex-col items-center min-h-svh py-8">
        <Health />
      </main>
    </Box>
  );
}
