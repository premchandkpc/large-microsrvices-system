import React, { useState, useEffect } from 'react';
import api from '../services/api';

export default function DashboardPage() {
  const [stats, setStats] = useState({
    totalDocuments: 0,
    processedDocuments: 0,
    pendingDocuments: 0,
    notifications: 0,
  });

  useEffect(() => {
    api.get('/documents?page=1&page_size=5').then((res) => {
      setStats((prev) => ({
        ...prev,
        totalDocuments: res.data?.total || 0,
      }));
    }).catch(() => {});
  }, []);

  return (
    <div>
      <div className="page-header">
        <h1>Dashboard</h1>
      </div>
      <div className="grid">
        <div className="card stat-card">
          <h3>{stats.totalDocuments}</h3>
          <p>Total Documents</p>
        </div>
        <div className="card stat-card">
          <h3>{stats.processedDocuments}</h3>
          <p>Processed</p>
        </div>
        <div className="card stat-card">
          <h3>{stats.pendingDocuments}</h3>
          <p>Pending</p>
        </div>
        <div className="card stat-card">
          <h3>{stats.notifications}</h3>
          <p>Notifications</p>
        </div>
      </div>
    </div>
  );
}
