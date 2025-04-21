import { useState } from "react";
import { LogOut } from "lucide-react";
import { useRouter } from "next/navigation";

export default function Navbar() {
	const router = useRouter();
	const [isDropdownOpen, setIsDropdownOpen] = useState(false);

	const toggleDropdown = () => {
		setIsDropdownOpen(!isDropdownOpen);
	};

	const submitLogout = async () => {
		await fetch('/api/logout')
		router.push("/login");
	};

	return (
		<div className="w-full">
			<div className="flex items-center justify-end pr-1">
			<div className="relative flex items-center justify-end py-1">
					<div
						onClick={toggleDropdown}
						className={`flex p-3 rounded-xl cursor-pointer font-semibold text-sm mr-3 transition ${isDropdownOpen ? "bg-white shadow-lg" : "hover:bg-white/20"}`}
					>
						{(() => {
							let name: string = ''
							if (typeof window !== 'undefined') {
								const user = localStorage.getItem('user')
								if (user) name = String(JSON.parse(user).username)
							}
							const username = name ?? 'Guest';
							return username.charAt(0).toUpperCase() + username.slice(1);
						})()}
					</div>
					{isDropdownOpen && (
						<>
						<div className="fixed inset-0 z-40" onClick={toggleDropdown}></div>
						<div className="absolute top-13 right-3 w-56 rounded-xl transition-all shadow-lg bg-white text-gray-700 p-2 z-50">
							<div className="space-y-3">
								<button
									onClick={submitLogout}
									className="flex items-center w-full px-3 py-2 rounded-lg hover:bg-gray-100 transition"
								>
									<LogOut className="mr-2 h-4 w-4" /> Sign Out
								</button>
							</div>
						</div>
						</>
					)}
				</div>
			</div>
		</div>
	);
}
