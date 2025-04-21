'use client'
import { LayoutDashboard, Users } from "lucide-react";
import { AppRouterInstance } from "next/dist/shared/lib/app-router-context.shared-runtime";
import { usePathname, useRouter  } from "next/navigation";

import React from "react";

interface SidebarMenuProps {
  icon: React.ComponentType<React.SVGAttributes<SVGSVGElement>>
  label: string; 
  active?: boolean;
  link: string;
  iconProps?: React.SVGProps<SVGSVGElement>;
}
export default function Sidebar() {
    const activePath = usePathname()
    const router = useRouter()

    const routes: SidebarMenuProps[] = [
      {
        icon : LayoutDashboard,
        label : 'Dashboard',
        active : false,
        iconProps : {className: "w-5 h-5"},
        link : '/dashboard'
      },
      {
        icon : Users,
        label : 'Customers',
        active : false,
        iconProps : {className: "w-5 h-5"},
        link : '/customers'
      }
    ]


    return (
      <div className="relative w-64 text-gray-800 flex flex-col p-4 shadow-xl mx-4 my-4 rounded-2xl bg-[#f4f4f5]">
        <div className="mb-5 text-gray-500 font-bold text-2xl">
          <span className="text-3xl">Superbank </span>
        </div>
        <nav className="flex-1 space-y-1">
          {routes.map((item, index) => (
            <SidebarLink key={index} router={router} label={item.label} link={item.link} icon={item.icon} iconProps={item.iconProps} active={activePath === item.link} />
          ))}
        </nav>
      </div>
    );
  }

  
  function SidebarLink({ label, active = false, icon: IconComponent, iconProps = {}, link, router}: SidebarMenuProps & {
    router : AppRouterInstance
  }) {

    return (
      <a
        onClick={() => {
          router.push(link)
        }}
        className={
          `cursor-pointer flex items-center space-x-3 rounded-md px-2 py-2 mb-3 text-sm font-medium transition-colors duration-150 gradient ${
          active ? "bg-gray-200 text-gray-900" : "hover:bg-gray-200 hover:text-gray-900 text-gray-700"
        }`}
      >

        <IconComponent {...iconProps} />
        <span>{label}</span>
      </a>
    );
  }
  
  