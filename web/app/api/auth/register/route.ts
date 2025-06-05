import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
const API_VERSION = process.env.API_VERSION || 'v1';

export async function POST(request: NextRequest) {
  try {
    const userData = await request.json();
    const { email, username, password, first_name, last_name } = userData;

    if (!email || !username || !password || !first_name || !last_name) {
      return NextResponse.json(
        { error: 'Tous les champs sont requis' },
        { status: 400 }
      );
    }

    // Appel à l'API Go pour l'inscription
    const apiResponse = await fetch(`${API_BASE_URL}/api/${API_VERSION}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(userData),
    });

    const data = await apiResponse.json();

    if (!apiResponse.ok) {
      return NextResponse.json(
        { error: data.message || 'Erreur lors de l\'inscription' },
        { status: apiResponse.status }
      );
    }

    // Stocker le token dans un cookie httpOnly
    const cookieStore = await cookies();
    cookieStore.set('auth-token', data.token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 60 * 60 * 24 * 7, // 7 jours
      path: '/',
    });

    // Retourner les données utilisateur sans le token
    return NextResponse.json({
      message: data.message,
      user: data.user,
    });

  } catch (error) {
    console.error('Erreur lors de l\'inscription:', error);
    return NextResponse.json(
      { error: 'Erreur serveur lors de l\'inscription' },
      { status: 500 }
    );
  }
} 