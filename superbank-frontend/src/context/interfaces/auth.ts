export interface LoginRequestBody {
  email: string;
  password: string;
}

export interface UserProfileResponse {
  id : string
  email : string
  username : string
}

export interface LoginResponseBody {
  user: UserProfileResponse
  token: string;
}