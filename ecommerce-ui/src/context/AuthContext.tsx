import { createContext, useState, useContext, useEffect, type ReactNode } from 'react';
import axios from 'axios';
axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

import { useNavigate } from 'react-router-dom';

interface User { id: number; name: string; email: string; }
interface AuthContextType {
  user: User | null;
  token: string | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
   const storedToken = localStorage.getItem('token');
   const storedUser = localStorage.getItem('user');
   if (storedToken && storedUser) {
    let parsedUser = null;
     try {
       parsedUser = JSON.parse(storedUser);
    } catch {
       // bad JSON – clear it to avoid future errors
       localStorage.removeItem('user');
       localStorage.removeItem('token');
     }
     if (parsedUser) {
       setToken(storedToken);
       setUser(parsedUser);
       axios.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`;
     }
   }
  }, []);

  const login = async (email: string, password: string) => {
    const { data } = await axios.post('/users/login', { email, password });
    const { token: jwt, user: userData } = data;
    localStorage.setItem('token', jwt);
    localStorage.setItem('user', JSON.stringify(userData));
    axios.defaults.headers.common['Authorization'] = `Bearer ${jwt}`;
    setToken(jwt);
    setUser(userData);
    navigate('/');
  };

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('cartID');  
    delete axios.defaults.headers.common['Authorization'];
    setToken(null);
    setUser(null);
    navigate('/login');
  };

  return (
    <AuthContext.Provider value={{ user, token, login, logout, isAuthenticated: !!token }}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used inside AuthProvider');
  return ctx;
};
