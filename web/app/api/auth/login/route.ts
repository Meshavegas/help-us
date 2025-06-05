import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
const API_VERSION = process.env.API_VERSION || 'v1';

export async function POST(request: NextRequest) {
  try {
    const { email, password } = await request.json();

    if (!email || !password) {
      return NextResponse.json(
        { error: 'Email et mot de passe requis' },
        { status: 400 }
      );
    }

    console.log('🔐 Tentative de connexion pour:', email);

    // Appel à l'API Go pour l'authentification
    const apiResponse = await fetch(`${API_BASE_URL}/api/${API_VERSION}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    const data = await apiResponse.json();

    if (!apiResponse.ok) {
      console.log('❌ Échec de connexion:', data.error || data.message);
      return NextResponse.json(
        { error: data.error || data.message || 'Erreur de connexion' },
        { status: apiResponse.status }
      );
    }

    console.log('✅ Connexion réussie pour:', email);
    console.log('🔑 Token reçu:', data.token ? `${data.token.substring(0, 20)}...` : 'Aucun');

    // Stocker le token dans un cookie httpOnly
    const cookieStore = await cookies();
    cookieStore.set('auth-token', data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 7, // 7 jours
      path: '/',
    });

    console.log('🍪 Cookie auth-token défini');

    // Retourner les données utilisateur sans le token
    return NextResponse.json({
      message: data.message || 'Connexion réussie',
      user: data.user,
    });

  } catch (error) {
    console.error('❌ Erreur lors de la connexion:', error);
    return NextResponse.json(
      { error: 'Erreur serveur lors de la connexion' },
      { status: 500 }
    );
  }
} 