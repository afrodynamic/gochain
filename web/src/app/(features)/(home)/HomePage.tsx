import { ChainDashboard } from '@/app/(features)/(chain)/ChainDashboard';
import { Navbar } from '@/app/_components/layout/Navbar';
import { Box } from '@mui/material';

export default function HomePage() {
  return (
    <Box className="min-h-svh">
      <Navbar />

      <main className="w-full">
        <ChainDashboard />
      </main>
    </Box>
  );
}
