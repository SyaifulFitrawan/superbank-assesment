import { cookies } from "next/headers";
import { NextResponse } from "next/server";

export async function GET() {
  const cookieStore = await cookies();
  cookieStore.delete("Authorization");
  return NextResponse.json({ message: "Logged out successfully" }, { status: 200 });
}