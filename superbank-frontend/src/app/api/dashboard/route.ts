import { routes } from "@/app/api";
import { BaseResponse } from "@/context/interfaces/base_response";
import { DashboardResponse } from "@/context/interfaces/dashboard";
import { cookies } from "next/headers";

export async function GET() {
  const ck = await cookies()
  const res = await fetch(`${routes.dashboard}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Authorization": String(ck.get('Authorization')?.value)
    },
  });

  const data: BaseResponse<DashboardResponse> = await res.json();

  return new Response(JSON.stringify(data), {
    headers: { 'Content-Type': 'application/json' },
  });
}
