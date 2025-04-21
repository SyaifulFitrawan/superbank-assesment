"use client";

import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import { CustomerData } from "@/context/interfaces/customer";
import moment from "moment";
import { BaseResponse } from "@/context/interfaces/base_response";

export default function CustomerDetailPage() {
  const router = useParams()
  const {id} = router

  const [customer, setCustomer] = useState<CustomerData>()

  const fetchCustomerDetail = async () => {
    const result: BaseResponse<CustomerData> = await fetch(`/api/customer/detail?id=${id}`).then(res => res.json());
    if (result && result.data) {
      setCustomer(result.data)
    }
  }

  useEffect(() => {
    fetchCustomerDetail()
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  let balanceWithPocket: number = customer?.pockets?.reduce((acc, curr) => acc + curr.balance, 0) ?? 0
  balanceWithPocket = (customer?.balance ?? 0) + balanceWithPocket
  return (
    <div className="px-3 mt-4 mb-4">
      <h2 className="font-bold text-lg text-gray-500">Customer Detail</h2>
      <p className="text-sm text-gray-500">Customers Detail. </p>
      <div className="flex mt-4 gap-6">
        <div className="w-1/3 shadow-lg rounded-xl p-4 space-y-4">
          <h2 className="text-xl text-gray-500 font-bold">Personal Info</h2>
          {[
            { label: "Username", value: customer ? customer.name : '-' },
            { label: "Parent Name", value: customer ? customer.parent_name : '-' },
            { label: "Phone", value: customer ? customer.phone : '-' },
            { label: "Account Number", value: customer ? customer.account_number : '-' },
            { label: "Account Branch", value: customer ? customer.account_branch : '-' },
            { label: "Account Type", value: customer ? customer.account_type : '-' },
            {
              label: "Account Balance",
              value: customer
                ? new Intl.NumberFormat('id-ID', {
                    style: 'currency',
                    currency: 'IDR',
                    minimumFractionDigits: 0,
                  }).format(customer.balance)
                : '-',
            },
            {
              label: "Balance With Pocket",
              value: customer
                ? new Intl.NumberFormat('id-ID', {
                    style: 'currency',
                    currency: 'IDR',
                    minimumFractionDigits: 0,
                  }).format(balanceWithPocket)
                : '-',
            },
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

          <div className="flex justify-between items-center pt-2">
            <p className="text-sm text-gray-600">Address</p>
            <div className="flex items-center space-x-2">
              <p className="text-sm font-medium text-gray-800">{customer ? customer.address : '-'}</p>
            </div>
          </div>
        </div>

        <div className="w-2/3">
          <div className="flex flex-col gap-6">
            <div className="w-full">
              <div className="shadow-lg rounded-xl p-4 space-y-4">
                <h2 className="text-xl text-gray-500 font-bold">Deposits</h2>
                {[
                  { label: "Total Deposits", value: customer ? customer.deposits?.length : '-' },
                  {
                    label: "Total Amount",
                    value: customer
                      ? new Intl.NumberFormat('id-ID', {
                          style: 'currency',
                          currency: 'IDR',
                          minimumFractionDigits: 0,
                        }).format(Number(customer.deposits?.reduce((acc, curr) => acc + curr.amount, 0)))
                      : '-',
                  },
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

                <div className="mt-4 h-35 overflow-auto rounded-b-lg">
                  <table className="min-w-full min-h-0 bg-white">
                    <thead className="bg-gray-200 text-gray-500">
                      <tr>
                        <th className="py-2 px-3 text-left rounded-tl-lg">Amount</th>
                        <th className="py-2 px-3 text-left">Interest Rate</th>
                        <th className="py-2 px-3 text-left">Term Months</th>
                        <th className="py-2 px-3 text-left">Start Date</th>
                        <th className="py-2 px-3 text-left">Maturity Date</th>
                        <th className="py-2 px-3 text-left rounded-tr-lg">Is Withdrawn</th>
                      </tr>
                    </thead>
                    <tbody className="text-gray-700 overflow-auto">
                      {customer?.deposits?.map((item, index) => (
                        <tr className="hover:bg-indigo-50 transition text-sm" key={index}>
                          <td className="py-2 px-3">
                            {
                              new Intl.NumberFormat('id-ID', {
                                style: 'currency',
                                currency: 'IDR',
                                minimumFractionDigits: 0,
                              }).format(item.amount)
                            }
                          </td>
                          <td className="py-2 px-3">{item.interest_rate.toFixed(2)}</td>
                          <td className="py-2 px-3">{item.term_months}</td>
                          <td className="py-2 px-3">{moment(item.start_date).format("YYYY-MM-DD")}</td>
                          <td className="py-2 px-3">{moment(item.maturity_date).format("YYYY-MM-DD")}</td>
                          <td className="py-2 px-3">{!item.is_withdrawn ? "No" : "Yes"}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
            <div className="w-full">
              <div className="shadow-lg rounded-xl p-4 space-y-4">
                <h2 className="text-xl text-gray-500 font-bold">Pockets</h2>
                {[
                  { label: "Total Pockets", value: customer ? customer.pockets?.length : '-' },
                  {
                    label: "Total Amount",
                    value: customer
                      ? new Intl.NumberFormat('id-ID', {
                          style: 'currency',
                          currency: 'IDR',
                          minimumFractionDigits: 0,
                        }).format(Number(customer.pockets?.reduce((acc, curr) => acc + curr.balance, 0)))
                      : '-',
                  },
                ].map((item, idx) => (
                  <div
                    key={idx}
                    className="flex justify-between items-center border-b  border-b-gray-500"
                  >
                    <p className="text-sm text-gray-600">{item.label}</p>
                    <div className="flex items-center space-x-2">
                      <p className="text-sm font-medium text-gray-800">
                        {item.value}
                      </p>
                    </div>
                  </div>
                ))}

                <div className="mt-4 h-35 overflow-auto rounded-b-lg">
                  <table className="min-w-full min-h-0 bg-white">
                    <thead className="bg-gray-200 text-gray-500">
                      <tr>
                        <th className="py-2 px-3 text-left rounded-tl-lg">Name</th>
                        <th className="py-2 px-3 text-left">Balance</th>
                        <th className="py-2 px-3 text-left">Target Amount</th>
                        <th className="py-2 px-3 text-left">Target Date</th>
                        <th className="py-2 px-3 text-left rounded-tr-lg">Is Active</th>
                      </tr>
                    </thead>
                    <tbody className="text-gray-700 overflow-auto">
                      {customer?.pockets?.map((item, index) => (
                        <tr className="hover:bg-indigo-50 transition text-sm" key={index}>
                          <td className="py-2 px-3">{item.name ?? '-'}</td>
                          <td className="py-2 px-3">
                            {
                              new Intl.NumberFormat('id-ID', {
                                style: 'currency',
                                currency: 'IDR',
                                minimumFractionDigits: 0,
                              }).format(item.balance)
                            }
                          </td>
                          <td className="py-2 px-3">{item.target_amount ?? '-'}</td>
                          <td className="py-2 px-3">{item.target ?? '-'}</td>
                          <td className="py-2 px-3">{!item.is_active ? 'No' : "Yes"}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
