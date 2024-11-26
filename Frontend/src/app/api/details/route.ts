import { NextRequest, NextResponse } from "next/server";

export async function GET(request: NextRequest) {
  try {
    if (process.env.NEXT_PUBLIC_BACKEND_URL) {
    const id = request.nextUrl.searchParams.get("id");
      const req = new NextRequest(
        `${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/public/points/${id}`,
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
