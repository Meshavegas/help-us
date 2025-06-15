import { getCurrentUser } from '@/lib/auth';
import { redirect } from 'next/navigation';

export default async function PaymentsPage() {
  const user = await getCurrentUser();
  if (!user) redirect('/login');

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Payments</h1>
      <p className="text-muted-foreground">Gestion des paiements prochainement disponible.</p>
    </div>
  );
} 