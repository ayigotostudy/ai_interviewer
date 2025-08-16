import instance from './axiosInstance'

export async function get(url: any, params: any) {
  const response = await instance.get(url, { params })
  return response.data
}

export async function post(url: any, data: any, config?: any) {
  const response = await instance.post(url, data, config)
  return response.data
}

export async function put(url: any, data: any) {
  const response = await instance.put(url, data)
  return response.data
}

export async function del(url: any, params: any) {
  const response = await instance.delete(url, { params })
  return response.data
}