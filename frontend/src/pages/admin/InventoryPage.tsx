import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import api from '../../services/api';
import { Book, PaginatedResponse } from '../../types';
import Pagination from '../../components/Pagination';

function InventoryPage() {
  const [pageIndex, setPageIndex] = useState(1);
  const [searchString, setSearchString] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['adminInventory', pageIndex, searchString],
    queryFn: async () => {
      const res = await api.get('/admin/inventory', {
        params: { pageIndex, pageSize: 10, searchString },
      });
      return res.data as PaginatedResponse<Book>;
    },
  });

  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({
    name: '', author: '', isbn: '', year: '', summary: '',
    price: '', quantity: '', bookTypeId: '', conditionId: '', genreId: '', publisherId: '',
  });

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/admin/inventory', {
        name: form.name, author: form.author, isbn: form.isbn,
        year: form.year ? Number(form.year) : null,
        summary: form.summary, price: Number(form.price), quantity: Number(form.quantity),
        bookTypeId: Number(form.bookTypeId), conditionId: Number(form.conditionId),
        genreId: Number(form.genreId), publisherId: Number(form.publisherId),
      });
      toast.success('Book added to inventory');
      setShowForm(false);
    } catch {
      toast.error('Failed to add book');
    }
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Inventory Management</h1>
      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
        <input placeholder="Search..." value={searchString}
          onChange={(e) => { setSearchString(e.target.value); setPageIndex(1); }}
          style={{ flex: 1, padding: '0.5rem' }} />
        <button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancel' : 'Add Book'}</button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} style={{ margin: '1rem 0', display: 'grid', gap: '0.5rem', maxWidth: '400px' }}>
          <input placeholder="Name" required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
          <input placeholder="Author" required value={form.author} onChange={(e) => setForm({ ...form, author: e.target.value })} />
          <input placeholder="ISBN" value={form.isbn} onChange={(e) => setForm({ ...form, isbn: e.target.value })} />
          <input placeholder="Year" type="number" value={form.year} onChange={(e) => setForm({ ...form, year: e.target.value })} />
          <input placeholder="Price" type="number" step="0.01" required value={form.price} onChange={(e) => setForm({ ...form, price: e.target.value })} />
          <input placeholder="Quantity" type="number" required value={form.quantity} onChange={(e) => setForm({ ...form, quantity: e.target.value })} />
          <input placeholder="Book Type ID" type="number" required value={form.bookTypeId} onChange={(e) => setForm({ ...form, bookTypeId: e.target.value })} />
          <input placeholder="Condition ID" type="number" required value={form.conditionId} onChange={(e) => setForm({ ...form, conditionId: e.target.value })} />
          <input placeholder="Genre ID" type="number" required value={form.genreId} onChange={(e) => setForm({ ...form, genreId: e.target.value })} />
          <input placeholder="Publisher ID" type="number" required value={form.publisherId} onChange={(e) => setForm({ ...form, publisherId: e.target.value })} />
          <textarea placeholder="Summary" value={form.summary} onChange={(e) => setForm({ ...form, summary: e.target.value })} />
          <button type="submit">Add Book</button>
        </form>
      )}

      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>ID</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Name</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Author</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Price</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Qty</th>
          </tr>
        </thead>
        <tbody>
          {data?.items.map((book) => (
            <tr key={book.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>{book.id}</td>
              <td style={{ padding: '0.5rem' }}>{book.name}</td>
              <td style={{ padding: '0.5rem' }}>{book.author}</td>
              <td style={{ padding: '0.5rem' }}>${book.price.toFixed(2)}</td>
              <td style={{ padding: '0.5rem' }}>{book.quantity}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
    </div>
  );
}

export default InventoryPage;
