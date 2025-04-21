import { routes } from "@/app/api";
import { BaseResponse } from "@/context/interfaces/base_response";
import { CustomerData } from "@/context/interfaces/customer";
import { cookies } from "next/headers";
import { NextRequest } from "next/server";

export async function GET(req: NextRequest) {
  const ck = await cookies()
  const {searchParams} = req.nextUrl

  console.log(searchParams.entries(), 'HAHA')
  const params = new URLSearchParams();
  searchParams.entries().forEach(([key, value]) => {
    if (value !== undefined && value !== null) {
      params.set(key, value.toString());
    }
  });

  const res = await fetch(`${routes.customer_list}?${params.toString()}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Authorization": String(ck.get('Authorization')?.value)
    },
  });

  const data: BaseResponse<CustomerData[]> = await res.json()

  return new Response(JSON.stringify(data), {
    headers: { 'Content-Type': 'application/json' },
  });
}