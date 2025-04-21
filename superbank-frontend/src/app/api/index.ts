const base_url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1` || "http://superbank-backend:8000/api/v1";

export const routes = {
  login: `${base_url}/login`,
  customer_list: `${base_url}/customer/list`,
  customer_detail: `${base_url}/customer/detail`,
  dashboard: `${base_url}/dashboard`
}