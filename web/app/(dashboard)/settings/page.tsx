import { getCurrentUser } from '@/lib/auth';
import { redirect } from 'next/navigation';

export default async function SettingsPage() {
  const user = await getCurrentUser();
  if (!user) redirect('/login');

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-4">Settings</h1>
      <p className="text-muted-foreground">Page de paramètres.</p>
    </div>
  );
} 