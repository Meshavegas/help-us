'use client';

export interface ApiResponse<T = any> {
  data?: T;
  message?: string;
  error?: string;
}

// Types pour l'API de la plateforme éducative
export interface User {
  id: number;
  username: string;
  email: string;
  phone_number?: string;
  role: 'famille' | 'enseignant' | 'administrator';
  is_active: boolean;
  created_at: string;
  updated_at: string;
  addresses?: Address[];
  payments?: Payment[];
  resources?: Resource[];
  famille?: Famille;
  enseignant?: Enseignant;
  administrator?: Administrator;
}

export interface Famille {
  user_id: number;
  family_name: string;
  missions?: Mission[];
  courses?: Course[];
  options?: Option[];
}

export interface Enseignant {
  user_id: number;
  specialization: string;
  qualifications: string;
  missions?: Mission[];
  courses?: Course[];
  reports?: Report[];
  options?: Option[];
  offers?: Offer[];
}

export interface Administrator {
  user_id: number;
  reports?: Report[];
  offers?: Offer[];
  resources?: Resource[];
}

export interface Address {
  id: number;
  user_id: number;
  street: string;
  city: string;
  postal_code: string;
  country: string;
  is_primary: boolean;
}

export interface Payment {
  id: number;
  user_id: number;
  amount: number;
  currency: string;
  status: string;
  payment_date: string;
  description?: string;
}

export interface Resource {
  id: number;
  title: string;
  description?: string;
  type: string;
  url?: string;
  content?: string;
  is_active: boolean;
}

export interface Mission {
  id: number;
  title: string;
  description?: string;
  status: string;
  start_date: string;
  end_date?: string;
}

export interface Course {
  id: number;
  title: string;
  description?: string;
  subject: string;
  level: string;
  duration: number;
  status: string;
}

export interface Report {
  id: number;
  title: string;
  content: string;
  status: string;
  created_at: string;
}

export interface Option {
  id: number;
  name: string;
  description?: string;
  price: number;
  is_active: boolean;
}

export interface Offer {
  id: number;
  title: string;
  description?: string;
  requirements?: string;
  salary?: number;
  is_active: boolean;
}

// Classe pour gérer les appels API côté client
class ApiClient {
  private baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      // Récupérer le token depuis le localStorage
      const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
      
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        headers: {
          'Content-Type': 'application/json',
          ...(token && { Authorization: `Bearer ${token}` }),
          ...options.headers,
        },
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

  // Méthodes d'authentification
  async login(email: string, password: string): Promise<ApiResponse<{ token: string; user: User }>> {
    const response = await this.request<{ token: string; user: User }>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });

    // Stocker le token si la connexion réussit
    if (response.data?.token && typeof window !== 'undefined') {
      localStorage.setItem('token', response.data.token);
    }

    return response;
  }

  async register(userData: {
    username: string;
    email: string;
    password: string;
    role: 'famille' | 'enseignant' | 'administrator';
    phone_number?: string;
    family_name?: string;
    specialization?: string;
    qualifications?: string;
  }): Promise<ApiResponse<{ token: string; user: User }>> {
    const response = await this.request<{ token: string; user: User }>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });

    // Stocker le token si l'inscription réussit
    if (response.data?.token && typeof window !== 'undefined') {
      localStorage.setItem('token', response.data.token);
    }

    return response;
  }

  async logout(): Promise<ApiResponse<void>> {
    const response = await this.request<void>('/api/v1/auth/logout', {
      method: 'POST',
    });

    // Supprimer le token du localStorage
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }

    return response;
  }

  async refreshToken(): Promise<ApiResponse<{ token: string }>> {
    return this.request<{ token: string }>('/api/v1/auth/refresh', {
      method: 'POST',
    });
  }

  // Méthodes pour la gestion des utilisateurs
  async getAllUsers(): Promise<ApiResponse<User[]>> {
    return this.request<User[]>('/api/v1/admin/users');
  }

  async getUser(id: number): Promise<ApiResponse<User>> {
    return this.request<User>(`/api/v1/users/${id}`);
  }

  async updateUser(id: number, userData: {
    username?: string;
    email?: string;
    phone_number?: string;
    family_name?: string;
    specialization?: string;
    qualifications?: string;
  }): Promise<ApiResponse<User>> {
    return this.request<User>(`/api/v1/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(userData),
    });
  }

  async deleteUser(id: number): Promise<ApiResponse<{ message: string }>> {
    return this.request<{ message: string }>(`/api/v1/admin/users/${id}`, {
      method: 'DELETE',
    });
  }

  // Méthodes pour les données associées aux utilisateurs
  async getUserAddresses(id: number): Promise<ApiResponse<Address[]>> {
    return this.request<Address[]>(`/api/v1/users/${id}/addresses`);
  }

  async getUserPayments(id: number): Promise<ApiResponse<Payment[]>> {
    return this.request<Payment[]>(`/api/v1/users/${id}/payments`);
  }

  async getUserResources(id: number): Promise<ApiResponse<Resource[]>> {
    return this.request<Resource[]>(`/api/v1/users/${id}/resources`);
  }

  // Méthodes pour le profil
  async getProfile(): Promise<ApiResponse<User>> {
    return this.request<User>('/api/v1/profile');
  }

  async updateProfile(userData: {
    username?: string;
    email?: string;
    phone_number?: string;
    family_name?: string;
    specialization?: string;
    qualifications?: string;
  }): Promise<ApiResponse<User>> {
    return this.request<User>('/api/v1/profile', {
      method: 'PUT',
      body: JSON.stringify(userData),
    });
  }

  // Méthodes génériques pour d'autres endpoints
  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint);
  }

  async post<T>(endpoint: string, data: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async put<T>(endpoint: string, data: any): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete<T>(endpoint: string): Promise<ApiResponse<T>> {
    return this.request<T>(endpoint, {
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