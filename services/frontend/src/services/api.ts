import axios from 'axios';

const API_BASE = process.env.REACT_APP_API_GATEWAY_URL || 'http://localhost:8081/api/v1';

const api = axios.create({
  baseURL: API_BASE,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;

// Auth
export const authAPI = {
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),
  register: (email: string, password: string, name: string) =>
    api.post('/auth/register', { email, password, name }),
  refresh: (refreshToken: string) =>
    api.post('/auth/refresh', { refresh_token: refreshToken }),
};

// Documents
export const documentsAPI = {
  list: (params?: any) => api.get('/documents', { params }),
  get: (id: string) => api.get(`/documents/${id}`),
  upload: (formData: FormData) =>
    api.post('/documents', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),
  process: (id: string, pipeline: string) =>
    api.post(`/documents/${id}/process`, { pipeline }),
};

// Search
export const searchAPI = {
  text: (query: string, filters?: string[]) =>
    api.get('/search', { params: { q: query, filter: filters } }),
  vector: (query: string, topK?: number) =>
    api.post('/search/vector', { query, top_k: topK }),
};

// Notifications
export const notificationsAPI = {
  list: () => api.get('/notifications'),
  markRead: (ids: string[]) =>
    api.post('/notifications/read', { notification_ids: ids }),
};

// Users
export const usersAPI = {
  getProfile: (id: string) => api.get(`/users/${id}`),
  updateProfile: (id: string, data: any) =>
    api.put(`/users/${id}/profile`, data),
};
