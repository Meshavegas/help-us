import { getCurrentUser } from '@/lib/auth';
import { redirect } from 'next/navigation';

export default async function LocationsPage() {
  const user = await getCurrentUser();
  if (!user) redirect('/login');

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Locations</h1>
      <p className="text-muted-foreground">Géolocalisation et adresses.</p>
    </div>
  );
} 