import React, { useState } from 'react';
import { searchAPI } from '../services/api';

export default function SearchPage() {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<any[]>([]);
  const [searching, setSearching] = useState(false);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) return;

    setSearching(true);
    try {
      const res = await searchAPI.text(query);
      setResults(res.data?.results || []);
    } catch (err) {
      console.error('Search failed', err);
    } finally {
      setSearching(false);
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1>Search Documents</h1>
      </div>

      <div className="card">
        <form onSubmit={handleSearch}>
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <input
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search documents..."
              style={{ flex: 1, padding: '0.75rem', borderRadius: 'var(--radius)', border: '1px solid var(--border)' }}
            />
            <button type="submit" className="btn btn-primary" disabled={searching}>
              {searching ? 'Searching...' : 'Search'}
            </button>
          </div>
        </form>
      </div>

      {results.length > 0 && (
        <div className="card">
          <h3 style={{ marginBottom: '1rem' }}>Results ({results.length})</h3>
          {results.map((result, i) => (
            <div key={i} style={{ padding: '1rem 0', borderBottom: '1px solid var(--border)' }}>
              <div style={{ marginBottom: '0.25rem' }}>
                <strong>{result.document_id}</strong>
                <span style={{ marginLeft: '0.5rem', color: 'var(--text-secondary)', fontSize: '0.85rem' }}>
                  Score: {result.score?.toFixed(3)}
                </span>
              </div>
              {result.snippet && (
                <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem' }}>
                  {result.snippet}
                </p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
