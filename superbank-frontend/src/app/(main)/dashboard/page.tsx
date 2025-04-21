"use client";

import { IPieChart } from "@/context/interfaces/chart";
import { DashboardResponse } from "@/context/interfaces/dashboard";
import { Users, BookLock, CreditCard } from "lucide-react";
import { useEffect, useState } from "react";
import { Pie } from "../_components/pie";
import { BaseResponse } from "@/context/interfaces/base_response";

export default function DashboardPage() {
  const [data, setData] = useState<DashboardResponse>()
  const [typePie, setTypePie] = useState<IPieChart[]>([])
  const [depositPie, setDepositPie] = useState<IPieChart[]>([])
  const [pocketPie, setPocketPie] = useState<IPieChart[]>([])

  const fetchDashboard = async () => {
    const result: BaseResponse<DashboardResponse> = await fetch('/api/dashboard').then(res => res.json());
    if (result && result.data) {
      setData(result.data)
      const type: IPieChart[] = result.data.type.map((item) => {return { name: item.account_type, value: item.count }})
      const deposit: IPieChart[] = result.data.deposits.map((item) => {return { name: item.range_label, value: item.count }})
      const pocket: IPieChart[] = result.data.pockets.map((item) => {return { name: item.range_label, value: item.count }})

      setTypePie(type ?? [])
      setDepositPie(deposit ?? [])
      setPocketPie(pocket ?? [])
    }
  }

  useEffect(() => {
    fetchDashboard()
  }, [])

  const typeHash: {[x: string]: IPieChart} = {}
  typePie.forEach((element) => (typeHash[element.name] = element))

  const depositHash: {[x: string]: IPieChart} = {}
  depositPie.forEach((element) => (depositHash[element.name] = element))

  const pocketHash: {[x: string]: IPieChart} = {}
  pocketPie.forEach((element) => (pocketHash[element.name] = element))

  return (
    <div className="px-3 mt-4">
      <h2 className="font-bold text-lg text-gray-500">Dashboard</h2>
      <p className="text-sm text-gray-500">Dashboard CMS. </p>
      <div className="h-full mt-4">
        <div className="flex flex-row space-x-2 ">
          <div className="flex flex-col w-full space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
              <div className="flex gap-6 bg-gray-200 w-full shadow-md backdrop-blur-sm px-4 py-4 rounded-xl">
                <div className="px-3 py-3 bg-gray-500 backdrop-blur-2xl w-15 rounded-xl">
                  <Users className="text-gray-500" size={35} color="white" />
                </div>
                <div className="flex flex-col mt-2">
                  <span className="text-gray-500 text-sm font-bold">
                    Customers
                  </span>
                  <span className="text-gray-500 text-4xl font-bold">
                    {data?.total.total_customers}
                  </span>
                </div>
              </div>
              <div className="flex gap-6 bg-gray-200 w-full shadow-md backdrop-blur-sm px-4 py-4 rounded-xl">
                <div className="px-3 py-3 bg-gray-500 backdrop-blur-2xl w-15 rounded-xl">
                  <BookLock className="text-white" size={35} color="white" />
                </div>
                <div className="flex flex-col mt-2">
                  <span className="text-gray-500 text-sm font-bold">
                    Total Deposits
                  </span>
                  <span className="text-gray-500 text-4xl font-bold">
                    {data?.total.total_deposits}
                  </span>
                </div>
              </div>
              <div className="flex gap-6 bg-gray-200 w-full shadow-md backdrop-blur-sm px-4 py-4 rounded-xl">
                <div className="px-3 py-3 bg-gray-500 backdrop-blur-2xl w-15 rounded-xl">
                  <CreditCard className="text-white" size={35} color="white" />
                </div>
                <div className="flex flex-col mt-2">
                  <span className="text-gray-500 text-sm font-bold">
                    Total Pocket
                  </span>
                  <span className="text-gray-500 text-4xl font-bold">
                    {data?.total.total_pockets}
                  </span>
                </div>
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-2">
              <div className="bg-gray-200 w-full shadow-md backdrop-blur-sm px-4 py-4 rounded-xl">
                <div className="flex justify-between">
                  <Pie data={typePie} />
                  <div className="w-2/4">
                    <h2 className="text-xl text-gray-500 font-bold">Account Type</h2>
                    {[
                      { label: "Silver", value: typeHash["Silver"]?.value ?? '-'},
                      { label: "Gold", value: typeHash["Gold"]?.value ?? '-'},
                      { label: "Platinum", value: typeHash["Platinum"]?.value ?? '-'},
                    ].map((item, idx) => (
                      <div
                        key={idx}
                        className="flex justify-between items-center py-2 border-b  border-b-gray-500"
                      >
                        <p className="text-sm text-gray-600">{item.label}</p>
                        <div className="flex items-center space-x-2">
                          <p className="text-sm font-medium text-gray-800">
                            {item.value}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              <div className="bg-gray-200 w-full shadow-md backdrop-blur-sm px-4 py-4 rounded-xl">
                <div className="flex justify-between">
                  <Pie data={depositPie} />
                  <div className="w-2/4">
                    <h2 className="text-xl text-gray-500 font-bold">Deposits</h2>
                    {[
                      { label: "0-1", value: depositHash["0-1"]?.value ?? '-'},
                      { label: "2-3", value: depositHash["2-3"]?.value ?? '-'},
                      { label: "4-5", value: depositHash["4-5"]?.value ?? '-'},
                      { label: "6+", value: depositHash["6+"]?.value ?? '-'},
                    ].map((item, idx) => (
                      <div
                        key={idx}
                        className="flex justify-between items-center py-2 border-b  border-b-gray-500"
                      >
                        <p className="text-sm text-gray-600">{item.label}</p>
                        <div className="flex items-center space-x-2">
                          <p className="text-sm font-medium text-gray-800">
                            {item.value}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
              <div className="bg-gray-200 w-full shadow-md backdrop-blur-sm px-4 py-4 rounded-xl">
                <div className="flex justify-between">
                  <Pie data={pocketPie} />
                  <div className="w-2/4">
                    <h2 className="text-xl text-gray-500 font-bold">Pockets</h2>
                    {[
                      { label: "0-1", value: pocketHash["0-1"]?.value ?? '-'},
                      { label: "2-3", value: pocketHash["2-3"]?.value ?? '-'},
                      { label: "4-5", value: pocketHash["4-5"]?.value ?? '-'},
                      { label: "6+", value: pocketHash["6+"]?.value ?? '-'},
                    ].map((item, idx) => (
                      <div
                        key={idx}
                        className="flex justify-between items-center py-2 border-b  border-b-gray-500"
                      >
                        <p className="text-sm text-gray-600">{item.label}</p>
                        <div className="flex items-center space-x-2">
                          <p className="text-sm font-medium text-gray-800">
                            {item.value}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
