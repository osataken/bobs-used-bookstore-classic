import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import api from '../services/api';
import { Offer, ReferenceDataItem, OfferStatusLabels } from '../types';

function ResalePage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({
    bookName: '', author: '', isbn: '', genreId: 0, conditionId: 0,
    publisherId: 0, bookTypeId: 0, summary: '', bookPrice: 0,
  });

  const { data: offers } = useQuery({
    queryKey: ['resaleOffers'],
    queryFn: async () => {
      const res = await api.get('/resale');
      return res.data.offers as Offer[];
    },
  });

  const { data: refData } = useQuery({
    queryKey: ['resaleRefData'],
    queryFn: async () => {
      const res = await api.get('/resale/referenceData');
      return res.data.referenceData as ReferenceDataItem[];
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/resale', form);
      toast.success('Offer submitted');
      setShowForm(false);
      queryClient.invalidateQueries({ queryKey: ['resaleOffers'] });
    } catch {
      toast.error('Failed to submit offer');
    }
  };

  const genres = refData?.filter((r) => r.dataType === 3) || [];
  const conditions = refData?.filter((r) => r.dataType === 1) || [];
  const publishers = refData?.filter((r) => r.dataType === 0) || [];
  const bookTypes = refData?.filter((r) => r.dataType === 2) || [];

  return (
    <div>
      <h1>Resale Offers</h1>
      <button onClick={() => setShowForm(!showForm)}>
        {showForm ? 'Cancel' : 'Submit New Offer'}
      </button>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ margin: '1rem 0', display: 'grid', gap: '0.5rem', maxWidth: '400px' }}>
          <input placeholder="Book Name" required value={form.bookName} onChange={(e) => setForm({ ...form, bookName: e.target.value })} />
          <input placeholder="Author" required value={form.author} onChange={(e) => setForm({ ...form, author: e.target.value })} />
          <input placeholder="ISBN" value={form.isbn} onChange={(e) => setForm({ ...form, isbn: e.target.value })} />
          <input type="number" step="0.01" placeholder="Price" required value={form.bookPrice || ''} onChange={(e) => setForm({ ...form, bookPrice: Number(e.target.value) })} />
          <select required value={form.genreId} onChange={(e) => setForm({ ...form, genreId: Number(e.target.value) })}>
            <option value={0}>Select Genre</option>
            {genres.map((g) => <option key={g.id} value={g.id}>{g.text}</option>)}
          </select>
          <select required value={form.conditionId} onChange={(e) => setForm({ ...form, conditionId: Number(e.target.value) })}>
            <option value={0}>Select Condition</option>
            {conditions.map((c) => <option key={c.id} value={c.id}>{c.text}</option>)}
          </select>
          <select required value={form.publisherId} onChange={(e) => setForm({ ...form, publisherId: Number(e.target.value) })}>
            <option value={0}>Select Publisher</option>
            {publishers.map((p) => <option key={p.id} value={p.id}>{p.text}</option>)}
          </select>
          <select required value={form.bookTypeId} onChange={(e) => setForm({ ...form, bookTypeId: Number(e.target.value) })}>
            <option value={0}>Select Book Type</option>
            {bookTypes.map((bt) => <option key={bt.id} value={bt.id}>{bt.text}</option>)}
          </select>
          <textarea placeholder="Summary" value={form.summary} onChange={(e) => setForm({ ...form, summary: e.target.value })} />
          <button type="submit">Submit Offer</button>
        </form>
      )}

      <h2>My Offers</h2>
      {!offers || offers.length === 0 ? (
        <p>No offers submitted yet.</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Book</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Author</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Price</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Status</th>
            </tr>
          </thead>
          <tbody>
            {offers.map((offer) => (
              <tr key={offer.id} style={{ borderBottom: '1px solid #ddd' }}>
                <td style={{ padding: '0.5rem' }}>{offer.bookName}</td>
                <td style={{ padding: '0.5rem' }}>{offer.author}</td>
                <td style={{ padding: '0.5rem' }}>${offer.bookPrice.toFixed(2)}</td>
                <td style={{ padding: '0.5rem' }}>{OfferStatusLabels[offer.offerStatus]}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default ResalePage;
