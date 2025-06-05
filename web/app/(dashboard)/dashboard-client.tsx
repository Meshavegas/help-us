'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient, useApi } from '@/lib/api';
import { User } from '@/lib/auth';

interface DashboardClientProps {
  user: User;
}

interface UsersResponse {
  users: User[];
  count: number;
}

export default function DashboardClient({ user }: DashboardClientProps) {
  const router = useRouter();
  const [users, setUsers] = useState<User[]>([]);
  const { loading, error, data, execute } = useApi<UsersResponse>();

  useEffect(() => {
    // Charger la liste des utilisateurs
    execute(() => apiClient.getUsers() as Promise<any>);
  }, [execute]);

  useEffect(() => {
    if (data?.users) {
      setUsers(data.users);
    }
  }, [data]);

  const handleLogout = async () => {
    try {
      await apiClient.logout();
      router.push('/login');
      router.refresh();
    } catch (error) {
      console.error('Erreur lors de la déconnexion:', error);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-6">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
              <p className="text-gray-600">
                Bienvenue, {user.first_name} {user.last_name}
              </p>
            </div>
            <button
              onClick={handleLogout}
              className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium"
            >
              Déconnexion
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          {/* User Profile Card */}
          <div className="bg-white overflow-hidden shadow rounded-lg mb-6">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900 mb-4">
                Profil Utilisateur
              </h3>
              <dl className="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
                <div>
                  <dt className="text-sm font-medium text-gray-500">Nom complet</dt>
                  <dd className="mt-1 text-sm text-gray-900">
                    {user.first_name} {user.last_name}
                  </dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Email</dt>
                  <dd className="mt-1 text-sm text-gray-900">{user.email}</dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Nom d'utilisateur</dt>
                  <dd className="mt-1 text-sm text-gray-900">{user.username}</dd>
                </div>
                <div>
                  <dt className="text-sm font-medium text-gray-500">Statut</dt>
                  <dd className="mt-1 text-sm text-gray-900">
                    <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                      user.is_active 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {user.is_active ? 'Actif' : 'Inactif'}
                    </span>
                  </dd>
                </div>
              </dl>
            </div>
          </div>

          {/* Users List */}
          <div className="bg-white shadow overflow-hidden sm:rounded-md">
            <div className="px-4 py-5 sm:px-6">
              <h3 className="text-lg leading-6 font-medium text-gray-900">
                Liste des Utilisateurs
              </h3>
              <p className="mt-1 max-w-2xl text-sm text-gray-500">
                Tous les utilisateurs enregistrés dans le système
              </p>
            </div>
            
            {loading && (
              <div className="px-4 py-5 sm:p-6">
                <div className="text-center">
                  <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
                  <p className="mt-2 text-sm text-gray-500">Chargement...</p>
                </div>
              </div>
            )}

            {error && (
              <div className="px-4 py-5 sm:p-6">
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
                  Erreur: {error}
                </div>
              </div>
            )}

            {!loading && !error && users.length > 0 && (
              <ul className="divide-y divide-gray-200">
                {users.map((userItem) => (
                  <li key={userItem.id} className="px-4 py-4 sm:px-6">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center">
                        <div className="flex-shrink-0">
                          <div className="h-10 w-10 rounded-full bg-indigo-500 flex items-center justify-center">
                            <span className="text-sm font-medium text-white">
                              {userItem.first_name[0]}{userItem.last_name[0]}
                            </span>
                          </div>
                        </div>
                        <div className="ml-4">
                          <div className="text-sm font-medium text-gray-900">
                            {userItem.first_name} {userItem.last_name}
                          </div>
                          <div className="text-sm text-gray-500">
                            {userItem.email} • @{userItem.username}
                          </div>
                        </div>
                      </div>
                      <div className="flex items-center">
                        <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                          userItem.is_active 
                            ? 'bg-green-100 text-green-800' 
                            : 'bg-red-100 text-red-800'
                        }`}>
                          {userItem.is_active ? 'Actif' : 'Inactif'}
                        </span>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            )}

            {!loading && !error && users.length === 0 && (
              <div className="px-4 py-5 sm:p-6">
                <p className="text-center text-gray-500">Aucun utilisateur trouvé</p>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  );
} 