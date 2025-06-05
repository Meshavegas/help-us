import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

export async function POST(request: NextRequest) {
  try {
    const cookieStore = await cookies();
    
    // Supprimer le cookie d'authentification
    cookieStore.delete('auth-token');
    
    return NextResponse.json(
      { message: 'Déconnexion réussie' },
      { status: 200 }
    );
  } catch (error) {
    console.error('Erreur lors de la déconnexion:', error);
    return NextResponse.json(
      { error: 'Erreur serveur lors de la déconnexion' },
      { status: 500 }
    );
  }
} 