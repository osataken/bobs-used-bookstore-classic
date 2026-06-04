import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { getAdminInventory } from '../../services/api';
import type { Book } from '../../types';

export default function AdminInventory() {
  const [pageIndex, setPageIndex] = useState(1);
  const [searchString, setSearchString] = useState('');

  const { data } = useQuery({
    queryKey: ['admin-inventory', searchString, pageIndex],
    queryFn: () => getAdminInventory({ searchString, pageIndex, pageSize: 10 }),
  });

  return (
    <div>
      <h1>Inventory Management</h1>
      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
        <input placeholder="Search..." value={searchString} onChange={e => { setSearchString(e.target.value); setPageIndex(1); }} style={{ flex: 1, padding: '0.5rem' }} />
        <Link to="/admin/inventory/new"><button>Add Book</button></Link>
      </div>

      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead><tr><th>Name</th><th>Author</th><th>Price</th><th>Qty</th><th>Actions</th></tr></thead>
        <tbody>
          {(data?.items as Book[])?.map(book => (
            <tr key={book.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>{book.name}</td>
              <td>{book.author}</td>
              <td>${book.price.toFixed(2)}</td>
              <td style={{ color: book.quantity === 0 ? 'red' : undefined }}>{book.quantity}</td>
              <td><Link to={`/admin/inventory/${book.id}/edit`}>Edit</Link></td>
            </tr>
          ))}
        </tbody>
      </table>

      {data && data.totalPages > 1 && (
        <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
          <button disabled={pageIndex <= 1} onClick={() => setPageIndex(p => p - 1)}>Previous</button>
          <span>Page {pageIndex} of {data.totalPages}</span>
          <button disabled={pageIndex >= data.totalPages} onClick={() => setPageIndex(p => p + 1)}>Next</button>
        </div>
      )}
    </div>
  );
}
