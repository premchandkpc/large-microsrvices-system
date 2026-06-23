import React from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export default function Layout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <div className="layout">
      <aside className="sidebar">
        <h2>Platform</h2>
        <nav>
          <NavLink to="/" end>Dashboard</NavLink>
          <NavLink to="/documents">Documents</NavLink>
          <NavLink to="/search">Search</NavLink>
          <NavLink to="/notifications">Notifications</NavLink>
          {user?.roles?.includes('ROLE_ADMIN') && (
            <NavLink to="/admin">Admin</NavLink>
          )}
        </nav>
        <div style={{ marginTop: 'auto' }}>
          <div style={{ padding: '0.75rem 1rem', fontSize: '0.85rem', color: 'var(--text-secondary)' }}>
            {user?.name}
          </div>
          <button className="btn btn-secondary" onClick={handleLogout} style={{ width: '100%' }}>
            Logout
          </button>
        </div>
      </aside>
      <main className="content">
        <Outlet />
      </main>
    </div>
  );
}
