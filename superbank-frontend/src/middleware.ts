import { NextRequest, NextResponse } from 'next/server';
import { jwtVerify } from 'jose';
import { cookies } from 'next/headers';

const SECRET_KEY = new TextEncoder().encode(process.env.JWT_SECRET || 'default_secret');

async function verify(token: string) {
  try {
    const { payload } = await jwtVerify(token, SECRET_KEY);
    return payload;
  } catch {
    return null;
  }
}

export async function middleware(req: NextRequest) {
  const ck = await cookies()
  const { pathname } = req.nextUrl;
  const token: string = String(ck.get('Authorization')?.value)

  if(!token) {
    if (pathname !== '/login') {
      const loginUrl = req.nextUrl.clone();
      loginUrl.pathname = '/login';
      return NextResponse.rewrite(loginUrl);
    } else {
      return NextResponse.next()
    }
  }

  const user = await verify(token)

  if (!user) {
    if (pathname !== '/login') {
      const loginUrl = req.nextUrl.clone();
      loginUrl.pathname = '/login';
      return NextResponse.redirect(loginUrl);
    } else {
      return NextResponse.next();
    }
  }

  if (pathname === '/') {
    const dashboardUrl = req.nextUrl.clone();
    dashboardUrl.pathname = '/dashboard';
    return NextResponse.redirect(dashboardUrl);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/', '/dashboard', '/((?!api|_next|static|favicon.ico).*)'],
};