"use client";

import { Eye, Search } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import moment from "moment";
import { CustomerData, IQueryParams } from "@/context/interfaces/customer";
import { BaseResponse } from "@/context/interfaces/base_response";

export default function CustomerPage() {
  const router = useRouter();

  const [customers, setCustomers] = useState<CustomerData[]>([]);
  const [search, setSearch] = useState<string>("");
  const [totalPages, setTotalPages] = useState<number>(1);
  const [currentPage, setCurrentPage] = useState<number>(0);

  const fetchCustomerList = async (req: IQueryParams) => {
    const api: string = '/api/customer/list'
    const params: string = `?page=${req.page}&limit=${req.limit}&search=${req.search}`
    const result: BaseResponse<CustomerData[]> = await fetch(`${api}${params}`).then(res => res.json());
    setCustomers(result.data);

    if (result.paginator) {
      const pageCount = result.paginator.pageCount || 1;
      setTotalPages(pageCount);
    }
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
  };

  useEffect(() => {
    fetchCustomerList({
      page: currentPage,
      limit: 10,
      search: search,
    });
  }, [search, currentPage]);

  return (
    <div className="">
      <div className="px-3 mt-4 mb-4">
        <div className="flex justify-between">
          <div>
            <h2 className="font-bold text-lg text-gray-500">Customer list</h2>
            <p className="text-sm text-gray-500">Customers Data Table.</p>
          </div>
          <div>
            <form className="mx-auto flex justify-end mb-4">
              <label htmlFor="default-search" className="sr-only">
                Search
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
                  <Search className="w-4 h-4 text-gray-500" />
                </div>
                <input
                  type="search"
                  value={search}
                  onChange={(e) => {
                    setCurrentPage(0);
                    setSearch(e.target.value);
                  }}
                  id="default-search"
                  className="block w-full p-3 ps-8 text-sm border border-gray-200 rounded-lg bg-white text-gray-800 focus:outline-none focus:border-[#3a5a40]"
                  placeholder="Search..."
                />
              </div>
            </form>
          </div>
        </div>

        <table className="min-w-full bg-white rounded-lg">
          <thead className="bg-gray-200 text-gray-500">
            <tr>
              <th className="py-3 px-6 text-left rounded-tl-lg">Name</th>
              <th className="py-3 px-6 text-left">Account Number</th>
              <th className="py-3 px-6 text-left">Account Branch</th>
              <th className="py-3 px-6 text-left">Phone Number</th>
              <th className="py-3 px-6 text-left">Created At</th>
              <th className="py-3 px-6 text-left rounded-tr-lg">Action</th>
            </tr>
          </thead>
          <tbody className="text-gray-700">
            {customers.map((item, index) => (
              <tr className="hover:bg-indigo-50 transition text-sm" key={index}>
                <td className="py-3 px-6">{item.name}</td>
                <td className="py-3 px-6">{item.account_number}</td>
                <td className="py-3 px-6">{item.account_branch}</td>
                <td className="py-3 px-6">{item.phone}</td>
                <td className="py-3 px-6">{moment(item.createdAt).format("YYYY-MM-DD")}</td>
                <td className="py-3 px-6">
                  <button
                    className="cursor-pointer px-2 py-1 text-white bg-gray-500 rounded-full"
                    onClick={() => router.push(`/customers/${item.id}`)}
                  >
                    <Eye />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {/* Pagination */}
        <div className="flex items-center space-x-2 mt-2 justify-end">
          <div className="px-3 py-1 text-sm text-gray-500 bg-white rounded-md shadow">
            <span>Current {currentPage + 1}</span>
          </div>

          <div className="px-3 py-1 text-sm text-gray-500 bg-white rounded-md shadow">
            <span>Total {totalPages}</span>
          </div>

          <button
            className="px-3 py-1 text-sm text-gray-500 bg-white rounded-md shadow hover:bg-gray-100 disabled:opacity-50"
            onClick={() => handlePageChange(currentPage - 1)}
            disabled={currentPage === 0}
          >
            Prev
          </button>

          <button
            className="px-3 py-1 text-sm text-gray-500 bg-white rounded-md shadow hover:bg-gray-100 disabled:opacity-50"
            onClick={() => handlePageChange(currentPage + 1)}
            disabled={currentPage >= totalPages - 1}
          >
            Next
          </button>
        </div>
      </div>
    </div>
  );
}