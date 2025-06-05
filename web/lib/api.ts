'use client';

export interface ApiResponse<T = any> {
  data?: T;
  message?: string;
  error?: string;
}

// Classe pour gérer les appels API côté client
class ApiClient {
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(endpoint, {
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
        credentials: 'include', // Important pour envoyer les cookies
        ...options,
      });

      const data = await response.json();

      if (!response.ok) {
        return {
          error: data.error || data.message || `HTTP Error: ${response.status}`,
        };
      }

      return { data, message: data.message };
    } catch (error) {
      return {
        error: error instanceof Error ? error.message : 'Network error',
      };
    }
  }

  // Méthodes d'authentification (via les routes Next.js)
  async login(email: string, password: string) {
    return this.request('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
  }

  async register(userData: {
    email: string;
    username: string;
    password: string;
    first_name: string;
    last_name: string;
  }) {
    return this.request('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async logout() {
    return this.request('/api/auth/logout', {
      method: 'POST',
    });
  }

  // Méthodes pour les appels API proxifiés
  async getUsers() {
    return this.request(`${process.env.NEXTAUTH_URL || 'http://localhost:3000'}/api/users`);
  }

  async getUser(id: number) {
    return this.request(`/api/users/${id}`);
  }

  async updateUser(id: number, userData: Partial<{
    first_name: string;
    last_name: string;
    email: string;
    username: string;
    is_active: boolean;
  }>) {
    return this.request(`/api/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(userData),
    });
  }

  async deleteUser(id: number) {
    return this.request(`/api/users/${id}`, {
      method: 'DELETE',
    });
  }

  // Méthodes pour le profil
  async getProfile() {
    return this.request('/api/profile');
  }

  // Méthodes génériques pour d'autres endpoints
  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(`/api${endpoint}`);
  }

  async post<T>(endpoint: string, data: any): Promise<ApiResponse<T>> {
    return this.request<T>(`/api${endpoint}`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put<T>(endpoint: string, data: any): Promise<ApiResponse<T>> {
    return this.request<T>(`/api${endpoint}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(`/api${endpoint}`, {
      method: 'DELETE',
    });
  }
}

// Instance singleton
export const apiClient = new ApiClient();

// Hook personnalisé pour les appels API avec gestion d'état
import { useState, useCallback } from 'react';

export function useApi<T = any>() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<T | null>(null);

  const execute = useCallback(async (apiCall: () => Promise<ApiResponse<T>>) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await apiCall();
      
      if (response.error) {
        setError(response.error);
        setData(null);
      } else {
        setData(response.data || null);
        setError(null);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
      setData(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const reset = useCallback(() => {
    setLoading(false);
    setError(null);
    setData(null);
  }, []);

  return {
    loading,
    error,
    data,
    execute,
    reset,
  };
} 