import request from '@/plugins/axiosInstance'

export interface LoginParams {
  email: string;
  password: string;
}

export interface RegisterParams {
  email: string;
  password: string;
  auth_code: string;
}

export function login(data: LoginParams) {
  return request.post('/user/login', data)
}

export function register(data: RegisterParams) {
  return request.post('/user/register', data)
}

export function getAuthCode(email: string) {
  return request.get('/authcode', { 
    params: { email } 
  })
}