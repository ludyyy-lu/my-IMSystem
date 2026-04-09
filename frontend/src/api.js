// In dev mode with Vite proxy, BASE_URL is empty (relative URLs go through proxy).
// In production, set VITE_API_URL to point to the api-gateway host.
const BASE_URL = import.meta.env.VITE_API_URL || '';

function getToken() {
  return localStorage.getItem('access_token');
}

async function request(method, path, body) {
  const headers = { 'Content-Type': 'application/json' };
  const token = getToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });
  const json = await res.json();
  if (!res.ok || json.code !== 200) {
    throw new Error(json.message || 'Request failed');
  }
  return json.data;
}

export const api = {
  register: (username, password) =>
    request('POST', '/api/auth/register', { username, password }),

  login: (username, password, device_id = 'web') =>
    request('POST', '/api/auth/login', { username, password, device_id }),

  getFriends: () => request('GET', '/api/friends'),
  getFriendRequests: () => request('GET', '/api/friends/requests'),
  getSentFriendRequests: () => request('GET', '/api/friends/requests/sent'),
  sendFriendRequest: (to_user_id, remark = '') =>
    request('POST', '/api/friends/request', { to_user_id, remark }),
  respondFriendRequest: (request_id, action) =>
    request('PUT', '/api/friends/request', { request_id, action }),

  searchUser: (keyword) => request('GET', `/api/users/search?keyword=${encodeURIComponent(keyword)}`),
  getMyProfile: () => request('GET', '/api/users/me'),
  getProfile: (id) => request('GET', `/api/users/${id}`),
  updateProfile: (data) => request('PUT', '/api/users/profile', data),

  getChatHistory: (peer_id, limit = 50) =>
    request('GET', `/api/messages/${peer_id}?limit=${limit}`),
  getUnreadCount: () => request('GET', '/api/messages/unread'),
};
