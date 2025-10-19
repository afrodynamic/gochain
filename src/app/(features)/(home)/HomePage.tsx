import { Health } from '@/app/(features)/(dashboard)/Health';
import { KeyGenerator } from '@/app/(features)/(dashboard)/KeyGenerator';
import { Navbar } from '@/app/_components/layout/Navbar';
import { Box } from '@mui/material';

export default function HomePage() {
  return (
    <Box className="flex flex-col min-h-svh ">
      <Navbar />

      <main className="flex flex-col items-center py-8 space-y-4">
        <Health />

        <KeyGenerator />
      </main>
    </Box>
  );
}
