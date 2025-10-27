import { BlockExplorer } from '@/app/(features)/(chain)/BlockExplorer';
import { Navbar } from '@/app/_components/layout/Navbar';
import { Box } from '@mui/material';

export default function BlocksPage() {
  return (
    <Box className="min-h-svh flex flex-col">
      <Navbar />

      <main className="flex-1">
        <BlockExplorer />
      </main>
    </Box>
  );
}
