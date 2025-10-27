import { WalletDashboard } from '@/app/(features)/(chain)/WalletDashboard';
import { Navbar } from '@/app/_components/layout/Navbar';

export default function WalletPage() {
  return (
    <div className="min-h-svh flex flex-col">
      <Navbar />

      <main className="flex-1">
        <WalletDashboard />
      </main>
    </div>
  );
}
