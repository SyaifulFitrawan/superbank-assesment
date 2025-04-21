"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import Loading from "../../_components/svg_loadings";
import { BaseResponse } from "@/context/interfaces/base_response";
import { LoginResponseBody } from "@/context/interfaces/auth";

export default function LoginPage() {
  const router = useRouter();

  const [loader, setLoader] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async (e : FormEvent<HTMLFormElement>
  ) => {
    e.preventDefault();
    setLoader(true)
    const formData = new FormData();
    formData.append('email', email);
    formData.append('password', password);

    try {
      const result: BaseResponse<LoginResponseBody> = await fetch('/api/login', {
        method: 'POST',
        body: formData,
      }).then(res => res.json());

      if (result) {
        localStorage.setItem('user', JSON.stringify(result.data.user))
      }

      router.push("/dashboard");
    } catch (err) {
      console.error('Login error:', err);
    }
  };

  return (
    <div className="min-h-screen bg-[#e4e3de] flex items-center justify-center">
      <div className="bg-gray-200 shadow-lg rounded-3xl flex w-full max-w-xl overflow-hidden">
        <div className="w-full p-10 flex flex-col justify-center">
          <div className="mb-6">
            <h1 className="text-4xl font-bold mt-2">
              Sign In<span className="text-blue-500">.</span>
            </h1>
          </div>

          <form className="space-y-4" onSubmit={(e : FormEvent<HTMLFormElement>) => handleSubmit(e)}>
            <input
              type="text"
              placeholder="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-gray-400"
            />
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Password"
              className="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-gray-400"
            />
            <div className="flex gap-4 mt-4">
              <button
                type="submit"
                className="flex-1 px-4 py-3 rounded-xl bg-gray-500 text-white hover:bg-gray-600 shadow-md"
              >
                {loader ?
                  <Loading/>
                  : <></>}
                Submit
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
