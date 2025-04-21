interface TotalCount {
  total_customers: string
  total_deposits: string
  total_pockets: string
}

interface AccountType {
  account_type: string
  count: number
}

interface Grouping {
  range_label: string
  count: number
}

export interface DashboardResponse {
  total: TotalCount
  type: AccountType[]
  deposits: Grouping[]
  pockets: Grouping[]
}