import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
const API_VERSION = process.env.API_VERSION || 'v1';

async function makeAuthenticatedRequest(endpoint: string, options: RequestInit = {}) {
  const cookieStore = await cookies();
  const token = cookieStore.get('auth-token')?.value;
  
  

  if (!token) {
    throw new Error('Token d\'authentification manquant');
  }

  const response = await fetch(`${API_BASE_URL}/api/${API_VERSION}${endpoint}`, {
    ...options,
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  return response;
}

export async function GET(request: NextRequest) {
  try {
    const apiResponse = await makeAuthenticatedRequest('/users');
    const data = await apiResponse.json();

    if (!apiResponse.ok) {
      return NextResponse.json(
        { error: data.message || 'Erreur lors de la récupération des utilisateurs' },
        { status: apiResponse.status }
      );
    }

    return NextResponse.json(data);

  } catch (error) {
    console.error('Erreur lors de la récupération des utilisateurs:', error);
    
    if (error instanceof Error && error.message === 'Token d\'authentification manquant') {
      return NextResponse.json(
        { error: error.message },
        { status: 401 }
      );
    }

    return NextResponse.json(
      { error: 'Erreur serveur' },
      { status: 500 }
    );
  }
} 