import { NextRequest, NextResponse } from 'next/server';
import { cookies } from 'next/headers';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';
const API_VERSION = process.env.API_VERSION || 'v1';

export async function GET(request: NextRequest) {
    console.log('🔍 Profile API: Request received');
    
  try {
    console.log('🔍 Profile API: Request received');
    
    const cookieStore = await cookies();
    const token = cookieStore.get('auth-token')?.value;

    console.log(`🍪 Profile API: Token present: ${!!token}`);

    if (!token) {
      console.log('❌ Profile API: No token found');
      return NextResponse.json(
        { error: 'Token d\'authentification manquant' },
        { status: 401 }
      );
    }
    const url = `${API_BASE_URL}/api/${API_VERSION}/profile`;

    console.log('🔗 Profile API: Calling Go API...',url);

    // Appel à l'API Go
    const apiResponse = await fetch(url, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
    });

    const data = await apiResponse.json();

    if (!apiResponse.ok) {
      console.log('❌ Profile API: Go API call failed:', apiResponse.status, data);
      return NextResponse.json(
        { error: data.error || data.message || 'Erreur lors de la récupération du profil' },
        { status: apiResponse.status }
      );
    }

    console.log('✅ Profile API: Success');
    return NextResponse.json(data);

  } catch (error) {
    console.error('❌ Profile API error:', error);
    return NextResponse.json(
      { error: 'Erreur serveur' },
      { status: 500 }
    );
  }
} 