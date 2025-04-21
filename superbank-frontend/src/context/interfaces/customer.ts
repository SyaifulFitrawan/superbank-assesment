export interface IQueryParams {
	search: string
    page: number
	limit: number
}

interface DepositData {
    amount: number
    interest_rate: number
    term_months: number
    start_date: string
    maturity_date: string
    is_withdrawn: boolean
    note: string
}

interface PocketData {
    name: string
    balance: number
    target_amount: number
    target: string
    is_active: boolean
}

export interface CustomerData {
    id: string
    name: string
    phone: string
    address: string
    parent_name: string
    account_number: string
    account_branch: string
    account_type: string
    balance: number
    createdAt: string
    deposits?: DepositData[]
    pockets?: PocketData[]
}