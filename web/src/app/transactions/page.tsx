import { TransactionsTable } from '@/app/(features)/(chain)/TransactionsTable';
import { Navbar } from '@/app/_components/layout/Navbar';

export default function TransactionsPage() {
  return (
    <div className="min-h-svh flex flex-col">
      <Navbar />

      <main className="flex-1">
        <TransactionsTable />
      </main>
    </div>
  );
}
