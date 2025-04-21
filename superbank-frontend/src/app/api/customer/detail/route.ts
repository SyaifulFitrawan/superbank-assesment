import { routes } from "@/app/api";
import { BaseResponse } from "@/context/interfaces/base_response";
import { CustomerData } from "@/context/interfaces/customer";
import { cookies } from "next/headers";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const ck = await cookies()
  const { searchParams } = req.nextUrl

  const res = await fetch(`${routes.customer_detail}/${searchParams.get('id')}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Authorization": String(ck.get('Authorization')?.value)
    },
  });

  const data: BaseResponse<CustomerData> = await res.json()

  return new Response(JSON.stringify(data), {
    headers: { 'Content-Type': 'application/json' },
  });
}