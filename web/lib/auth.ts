import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
const API_VERSION = process.env.API_VERSION || 'v1';
export interface User {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  message: string;
  user: User;
  token: string;
}

// Fonction pour faire des appels API avec gestion des erreurs
async function apiCall(endpoint: string, options: RequestInit = {}) {
  const url = `${API_BASE_URL}/api/${API_VERSION}${endpoint}`;
  
  const response = await fetch(url, {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.text();
    throw new Error(`API Error: ${response.status} - ${error}`);
  }

  return response.json();
}

// Connexion utilisateur
export async function signIn(email: string, password: string): Promise<AuthResponse> {
  const response = await apiCall('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });

  // Stocker le token dans un cookie httpOnly
  const cookieStore = await cookies();
  cookieStore.set('auth-token', response.token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: 60 * 60 * 24 * 7, // 7 jours
    path: '/',
  });

  return response;
}

// Inscription utilisateur
export async function signUp(userData: {
  email: string;
  username: string;
  password: string;
  first_name: string;
  last_name: string;
}): Promise<AuthResponse> {
  const response = await apiCall('/auth/register', {
    method: 'POST',
    body: JSON.stringify(userData),
  });

  // Stocker le token dans un cookie httpOnly
  const cookieStore = await cookies();
  cookieStore.set('auth-token', response.token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    maxAge: 60 * 60 * 24 * 7, // 7 jours
    path: '/',
  });

  return response;
}

// Déconnexion
export async function signOut() {
  const cookieStore = await cookies();
  cookieStore.delete('auth-token');
}

// Récupérer l'utilisateur actuel (côté serveur)
export async function getCurrentUser(): Promise<User | null> {
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get('auth-token')?.value;

    if (!token) {
      console.log('🚫 getCurrentUser: No token found');
      return null;
    }

    console.log('🔍 getCurrentUser: Token found, calling Go API directly');

    // Appeler directement l'API Go avec le token
    const response = await apiCall('/profile', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    console.log('✅ getCurrentUser: Profile retrieved successfully');
    
    return response.user || response;
  } catch (error) {
    console.error('❌ getCurrentUser error:', error);
    return null;
  }
}

// Middleware d'authentification
export async function auth(request: NextRequest) {
  // Vérifier si request est défini
  if (!request || !request.nextUrl) {
    return NextResponse.next();
  }

  const { pathname } = request.nextUrl;

  // Debug: Log de la route et des cookies
  console.log(`🔍 Auth middleware: ${pathname}`);

  // Pages publiques qui ne nécessitent pas d'authentification
  const publicPaths = ['/login', '/register', '/', '/api'];
  const isPublicPath = publicPaths.some(path => 
    pathname === path || pathname.startsWith(path + '/')
  );

  // Si c'est une page publique, laisser passer sans vérification
  if (isPublicPath) {
    console.log(`✅ Public path: ${pathname}`);
    return NextResponse.next();
  }

  // Pour les pages protégées, vérifier le token
  try {
    const cookieStore = await cookies();
    const token = cookieStore.get('auth-token')?.value;

    console.log(`🍪 Token present for ${pathname}: ${!!token}`);
    if (token) {
      console.log(`🔑 Token preview: ${token.substring(0, 20)}...`);
    }

    // Si pas de token, rediriger vers login
    if (!token) {
      console.log(`🚫 No token, redirecting to login from: ${pathname}`);
      const loginUrl = new URL('/login', request.url);
      return NextResponse.redirect(loginUrl);
    }

    // Si token présent, laisser passer
    console.log(`✅ Authenticated access to: ${pathname}`);
    return NextResponse.next();

  } catch (error) {
    console.error('❌ Auth middleware error:', error);
    // En cas d'erreur, laisser passer pour éviter les boucles
    return NextResponse.next();
  }
}

// Fonction utilitaire pour faire des appels API authentifiés côté serveur
export async function authenticatedApiCall(endpoint: string, options: RequestInit = {}) {
  const cookieStore = await cookies();
  const token = cookieStore.get('auth-token')?.value;

  if (!token) {
    throw new Error('No authentication token found');
  }

  return apiCall(endpoint, {
    ...options,
    headers: {
      Authorization: `Bearer ${token}`,
      ...options.headers,
    },
  });
}
