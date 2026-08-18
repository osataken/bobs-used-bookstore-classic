import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { adminService } from '../../services';
import { Book } from '../../types';
import { useAuth } from '../../context/AuthContext';

export default function Inventory() {
  const { user } = useAuth();
  const [pageIndex, setPageIndex] = useState(1);
  const [search, setSearch] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['admin-inventory', pageIndex, search],
    queryFn: () => adminService.listBooks({ pageIndex, pageSize: 10, searchString: search }).then(r => r.data),
    enabled: !!user?.isAdmin,
  });

  if (!user?.isAdmin) return <p>Access denied.</p>;

  return (
    <div>
      <h1 className="page-title">Inventory Management</h1>
      <input placeholder="Search..." value={search} onChange={e => { setSearch(e.target.value); setPageIndex(1); }} style={{ marginBottom: '1rem' }} />
      {isLoading ? <p>Loading...</p> : (
        <>
          <table>
            <thead><tr><th>ID</th><th>Name</th><th>Author</th><th>Price</th><th>Qty</th><th>Genre</th></tr></thead>
            <tbody>
              {(data?.items as Book[])?.map(b => (
                <tr key={b.id}><td>{b.id}</td><td>{b.name}</td><td>{b.author}</td><td>${b.price.toFixed(2)}</td><td>{b.quantity}</td><td>{b.genre?.text}</td></tr>
              ))}
            </tbody>
          </table>
          {data && data.totalPages > 1 && (
            <div className="pagination">
              <button className="btn" disabled={pageIndex <= 1} onClick={() => setPageIndex(p => p - 1)}>Prev</button>
              <span>Page {data.pageIndex} of {data.totalPages}</span>
              <button className="btn" disabled={pageIndex >= data.totalPages} onClick={() => setPageIndex(p => p + 1)}>Next</button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
