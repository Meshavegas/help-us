import { getCurrentUser } from '@/lib/auth';
import { redirect } from 'next/navigation';

export default async function TeachersPage() {
  const user = await getCurrentUser();
  if (!user) redirect('/login');

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Teachers</h1>
      <p className="text-muted-foreground">Cette page est en cours de développement.</p>
    </div>
  );
} 