import { NextRequest } from "next/server";
import {cookies} from 'next/headers'
import { LoginRequestBody, LoginResponseBody } from "@/context/interfaces/auth";
import { routes } from "@/app/api";
import { BaseResponse } from "@/context/interfaces/base_response";

export async function POST(request: NextRequest) {
  const payload = await request.formData()
  const email = payload.get('email')
  const password =  payload.get('password')

  const body: LoginRequestBody = {
    email: String(email),
    password: String(password)
  }

  const res = await fetch(`${routes.login}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });

  const data : BaseResponse<LoginResponseBody> = await res.json();

  const ck = await cookies()
  ck.set('Authorization', data.data.token.split(' ')[1])

  return new Response(JSON.stringify(data), {status: 200})
}