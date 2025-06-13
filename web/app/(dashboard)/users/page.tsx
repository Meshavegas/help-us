import { redirect } from 'next/navigation';
import UserManagement from '@/components/user-management';
import { User } from '@/lib/api';

// Mock function to get user data - replace with actual auth implementation
async function getCurrentUser(): Promise<User | null> {
  // This should be replaced with your actual authentication logic
  // For now, returning a mock admin user
  return {
    id: 1,
    username: 'admin',
    email: 'admin@example.com',
    role: 'administrator',
    is_active: true,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  };
}

export default async function UsersPage() {
  const user = await getCurrentUser();

  // Redirect if not authenticated
  if (!user) {
    redirect('/login');
  }

  // Redirect if not administrator
  if (user.role !== 'administrator') {
    redirect('/dashboard');
  }

  return (
    <div className="container mx-auto py-6">
      <UserManagement currentUser={user} />
    </div>
  );
} 