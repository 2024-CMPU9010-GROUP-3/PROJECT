import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  try {
    if (process.env.NEXT_PUBLIC_BACKEND_URL) {
      const userid = request.nextUrl.searchParams.get("userid");
      const req = new NextRequest(
        `${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/history/${userid}?limit=10&offset=0`,
        request
      );
      const resp = await fetch(req);
      const j = await resp.json();
      return NextResponse.json(j);
    }
    return NextResponse.json({}, { status: 500 });
  } catch (error) {
    console.error("Error fetching point by id:", error);
    return NextResponse.json({}, { status: 500 });
  }
}

export async function POST(request: NextRequest) {
  if (process.env.NEXT_PUBLIC_BACKEND_URL) {
    const userid = request.nextUrl.searchParams.get("userid");
    const req = new NextRequest(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/history/${userid}`,
      request
    );
    const resp = await fetch(req);
    const j = await resp.json();
    return NextResponse.json(j);
  }
  return NextResponse.json([], { status: 500 });
}

export async function DELETE(request: NextRequest) {
  if (process.env.NEXT_PUBLIC_BACKEND_URL) {
    const userid = request.nextUrl.searchParams.get("userid");
    const req = new NextRequest(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/history/${userid}`,
      request
    );
    const resp = await fetch(req);
    const j = await resp.json();
    return NextResponse.json(j);
  }
  return NextResponse.json([], { status: 500 });
}
