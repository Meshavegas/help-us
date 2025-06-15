import { getCurrentUser } from '@/lib/auth';
import { redirect } from 'next/navigation';

export default async function MyPaymentsPage() {
  const user = await getCurrentUser();
  if (!user) redirect('/login');

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">My Payments</h1>
      <p className="text-muted-foreground">Historique de vos paiements.</p>
    </div>
  );
} 