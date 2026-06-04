import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getBook, createBook, updateBook, getReferenceData } from '../../services/api';
import { ReferenceDataType } from '../../types';

export default function AdminInventoryForm() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const isEdit = !!id;
  const { data: refData } = useQuery({ queryKey: ['reference-data'], queryFn: getReferenceData });

  const [form, setForm] = useState({
    name: '', author: '', isbn: '', publisherId: 0, bookTypeId: 0, genreId: 0, conditionId: 0,
    price: 0, quantity: 0, year: null as number | null, summary: '',
  });

  useEffect(() => {
    if (isEdit) {
      getBook(Number(id)).then(book => {
        setForm({
          name: book.name, author: book.author, isbn: book.isbn, publisherId: book.publisherId,
          bookTypeId: book.bookTypeId, genreId: book.genreId, conditionId: book.conditionId,
          price: book.price, quantity: book.quantity, year: book.year, summary: book.summary || '',
        });
      });
    }
  }, [id, isEdit]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isEdit) {
      await updateBook(Number(id), form);
      toast.success('Book updated');
    } else {
      await createBook(form);
      toast.success('Book added to inventory');
    }
    navigate('/admin/inventory');
  };

  const bookTypes = refData?.filter(r => r.dataType === ReferenceDataType.BookType) || [];
  const conditions = refData?.filter(r => r.dataType === ReferenceDataType.Condition) || [];
  const genres = refData?.filter(r => r.dataType === ReferenceDataType.Genre) || [];
  const publishers = refData?.filter(r => r.dataType === ReferenceDataType.Publisher) || [];

  return (
    <div>
      <h1>{isEdit ? 'Edit Book' : 'Add Book'}</h1>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: '1rem', maxWidth: '500px' }}>
        <input required placeholder="Name" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
        <input required placeholder="Author" value={form.author} onChange={e => setForm({ ...form, author: e.target.value })} />
        <input placeholder="ISBN" value={form.isbn} onChange={e => setForm({ ...form, isbn: e.target.value })} />
        <select required value={form.bookTypeId} onChange={e => setForm({ ...form, bookTypeId: Number(e.target.value) })}>
          <option value={0}>Select Book Type</option>
          {bookTypes.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}
        </select>
        <select required value={form.conditionId} onChange={e => setForm({ ...form, conditionId: Number(e.target.value) })}>
          <option value={0}>Select Condition</option>
          {conditions.map(c => <option key={c.id} value={c.id}>{c.text}</option>)}
        </select>
        <select required value={form.genreId} onChange={e => setForm({ ...form, genreId: Number(e.target.value) })}>
          <option value={0}>Select Genre</option>
          {genres.map(g => <option key={g.id} value={g.id}>{g.text}</option>)}
        </select>
        <select required value={form.publisherId} onChange={e => setForm({ ...form, publisherId: Number(e.target.value) })}>
          <option value={0}>Select Publisher</option>
          {publishers.map(p => <option key={p.id} value={p.id}>{p.text}</option>)}
        </select>
        <input required type="number" step="0.01" placeholder="Price" value={form.price || ''} onChange={e => setForm({ ...form, price: parseFloat(e.target.value) || 0 })} />
        <input required type="number" placeholder="Quantity" value={form.quantity || ''} onChange={e => setForm({ ...form, quantity: parseInt(e.target.value) || 0 })} />
        <input type="number" placeholder="Year (optional)" value={form.year || ''} onChange={e => setForm({ ...form, year: e.target.value ? parseInt(e.target.value) : null })} />
        <textarea placeholder="Summary" value={form.summary} onChange={e => setForm({ ...form, summary: e.target.value })} />
        <button type="submit">{isEdit ? 'Update' : 'Create'}</button>
      </form>
    </div>
  );
}
