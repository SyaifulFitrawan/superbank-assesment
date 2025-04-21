'use client';
import "@/app/globals.css";
import Sidebar from "./_components/sidebar";
import Navbar from "./_components/navbar";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <section className="flex w-full h-full bg-[#e4e3de] relative">
      <Sidebar />
      <div className="flex-col flex w-full">
      <Navbar />
        <div className="overflow-hidden mb-4 w-full h-full no-scrollbar pr-4">
          <div className="flex flex-col bg-[#f4f4f5] w-full h-full bg-blur p-3 rounded-xl shadow-xl">
            {children}
          </div>
        </div>
      </div>
    </section>
  );
}
