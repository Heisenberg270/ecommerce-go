import axios from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
});

// attach JWT if present
api.interceptors.request.use(config => {
  const token = localStorage.getItem("jwt");
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;

// only a type, not a runtime import
export type Product = {
  id: number;
  name: string;
  description?: string;
  price: number;
};
