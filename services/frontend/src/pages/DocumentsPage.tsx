import React, { useState, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { documentsAPI } from '../services/api';

export default function DocumentsPage() {
  const [documents, setDocuments] = useState<any[]>([]);
  const [uploading, setUploading] = useState(false);
  const fileInput = useRef<HTMLInputElement>(null);
  const navigate = useNavigate();

  const handleUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    setUploading(true);
    const formData = new FormData();
    formData.append('file', file);

    try {
      const res = await documentsAPI.upload(formData);
      setDocuments((prev) => [res.data, ...prev]);
    } catch (err) {
      console.error('Upload failed', err);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Documents</h1>
        <button
          className="btn btn-primary"
          onClick={() => fileInput.current?.click()}
          disabled={uploading}
        >
          {uploading ? 'Uploading...' : 'Upload Document'}
        </button>
        <input
          ref={fileInput}
          type="file"
          style={{ display: 'none' }}
          onChange={handleUpload}
        />
      </div>

      <div className="card">
        <table className="table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>Size</th>
              <th>Status</th>
              <th>Date</th>
            </tr>
          </thead>
          <tbody>
            {documents.map((doc) => (
              <tr
                key={doc.id}
                onClick={() => navigate(`/documents/${doc.id}`)}
                style={{ cursor: 'pointer' }}
              >
                <td>{doc.filename}</td>
                <td>{doc.content_type}</td>
                <td>{(doc.size_bytes / 1024).toFixed(1)} KB</td>
                <td>
                  <span className={`badge badge-${doc.status === 'completed' ? 'success' : 'warning'}`}>
                    {doc.status}
                  </span>
                </td>
                <td>{new Date(doc.created_at).toLocaleDateString()}</td>
              </tr>
            ))}
            {documents.length === 0 && (
              <tr>
                <td colSpan={5} style={{ textAlign: 'center', padding: '2rem', color: 'var(--text-secondary)' }}>
                  No documents yet. Upload your first document.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
