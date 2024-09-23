import type { NextApiRequest, NextApiResponse } from 'next';

type Location = {
  lat: number;
  lng: number;
  name: string;
};

export default async function handler(req: NextApiRequest, res: NextApiResponse<Location[]>) {
  try {
    const response = await fetch('http://localhost:8000/api/locations');
    const data: Location[] = await response.json();
    res.status(200).json(data);
  } catch (error) {
    console.error('Error fetching locations:', error);
    res.status(500).json([]);
  }
}