import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { createOffer, getReferenceData } from '../services/api';
import { ReferenceDataType } from '../types';

export default function ResaleCreatePage() {
  const navigate = useNavigate();
  const { data: refData } = useQuery({ queryKey: ['reference-data'], queryFn: getReferenceData });
  const [form, setForm] = useState({ bookName: '', author: '', isbn: '', bookTypeId: 0, conditionId: 0, genreId: 0, publisherId: 0, bookPrice: 0 });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await createOffer(form);
    toast.success('Offer submitted');
    navigate('/resale');
  };

  const bookTypes = refData?.filter(r => r.dataType === ReferenceDataType.BookType) || [];
  const conditions = refData?.filter(r => r.dataType === ReferenceDataType.Condition) || [];
  const genres = refData?.filter(r => r.dataType === ReferenceDataType.Genre) || [];
  const publishers = refData?.filter(r => r.dataType === ReferenceDataType.Publisher) || [];

  return (
    <div>
      <h1>Submit Resale Offer</h1>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: '1rem', maxWidth: '500px' }}>
        <input required placeholder="Book Name" value={form.bookName} onChange={e => setForm({ ...form, bookName: e.target.value })} />
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
        <input required type="number" step="0.01" placeholder="Price" value={form.bookPrice || ''} onChange={e => setForm({ ...form, bookPrice: parseFloat(e.target.value) || 0 })} />
        <button type="submit">Submit Offer</button>
      </form>
    </div>
  );
}
