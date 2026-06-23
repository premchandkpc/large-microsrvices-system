import React, { useState, useEffect } from 'react';
import { notificationsAPI } from '../services/api';

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<any[]>([]);

  useEffect(() => {
    notificationsAPI.list().then((res) => {
      setNotifications(res.data || []);
    }).catch(() => {});
  }, []);

  const handleMarkRead = async (id: string) => {
    try {
      await notificationsAPI.markRead([id]);
      setNotifications((prev) =>
        prev.map((n) => (n.id === id ? { ...n, is_read: true } : n))
      );
    } catch (err) {
      console.error('Failed to mark as read', err);
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Notifications</h1>
      </div>
      <div className="card">
        {notifications.length === 0 ? (
          <p style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-secondary)' }}>
            No notifications yet
          </p>
        ) : (
          notifications.map((notif) => (
            <div
              key={notif.id}
              style={{
                padding: '1rem',
                borderBottom: '1px solid var(--border)',
                background: notif.is_read ? 'transparent' : '#f0f7ff',
                cursor: 'pointer',
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
              }}
              onClick={() => !notif.is_read && handleMarkRead(notif.id)}
            >
              <div>
                <strong>{notif.title}</strong>
                <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>
                  {notif.body}
                </p>
                <small style={{ color: 'var(--text-secondary)' }}>
                  {new Date(notif.created_at).toLocaleString()}
                </small>
              </div>
              {!notif.is_read && (
                <span className="badge badge-success">New</span>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
}
