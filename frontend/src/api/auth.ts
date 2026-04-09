import client from './client'

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface AuthResponse {
  token: string
  user: {
    id: string
    email: string
    name: string
    role: string
  }
}

export const login = (data: LoginRequest) =>
  client.post<AuthResponse>('/auth/login', data)

export const register = (data: RegisterRequest) =>
  client.post<AuthResponse>('/auth/register', data)

export const forgotPassword = (email: string) =>
  client.post('/auth/forgot-password', { email })

export const resetPassword = (token: string, password: string) =>
  client.post('/auth/reset-password', { token, password })
