import { NextResponse } from 'next/server';

export async function GET() {
  try {
    const mapBoxApiKey = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN;
    if (!mapBoxApiKey) {
      throw new Error('Mapbox API key not found');
    }

    return NextResponse.json({ mapBoxApiKey });
  } catch (error) {
    console.error('Error fetching Mapbox API key:', error);
    return NextResponse.json({ mapBoxApiKey: '' }, { status: 500 });
  }
}
