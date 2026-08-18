import { useState, useEffect } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { resaleService } from '../services';
import { ReferenceDataItem, OfferStatusLabels } from '../types';
import { useAuth } from '../context/AuthContext';
import toast from 'react-hot-toast';

export default function Resale() {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const { data: offers, isLoading } = useQuery({
    queryKey: ['resale'],
    queryFn: () => resaleService.list().then(r => r.data),
    enabled: !!user,
  });

  const [showForm, setShowForm] = useState(false);
  const [refData, setRefData] = useState<{ publishers: ReferenceDataItem[]; conditions: ReferenceDataItem[]; bookTypes: ReferenceDataItem[]; genres: ReferenceDataItem[] }>({ publishers: [], conditions: [], bookTypes: [], genres: [] });
  const [form, setForm] = useState({ bookName: '', author: '', isbn: '', bookTypeId: 0, conditionId: 0, genreId: 0, publisherId: 0, bookPrice: 0, summary: '' });

  useEffect(() => {
    if (showForm) {
      resaleService.getReferenceData().then(r => setRefData(r.data));
    }
  }, [showForm]);

  const submit = async () => {
    try {
      await resaleService.create(form);
      queryClient.invalidateQueries({ queryKey: ['resale'] });
      toast.success('Offer submitted');
      setShowForm(false);
    } catch { toast.error('Failed to submit'); }
  };

  if (!user) return <p>Please log in.</p>;
  if (isLoading) return <p>Loading...</p>;

  return (
    <div>
      <h1 className="page-title">Sell Your Books</h1>
      <button className="btn btn-primary" onClick={() => setShowForm(true)} style={{ marginBottom: '1rem' }}>Submit Offer</button>
      {showForm && (
        <div className="card">
          <div className="form-group"><label>Book Name</label><input value={form.bookName} onChange={e => setForm({ ...form, bookName: e.target.value })} /></div>
          <div className="form-group"><label>Author</label><input value={form.author} onChange={e => setForm({ ...form, author: e.target.value })} /></div>
          <div className="form-group"><label>ISBN</label><input value={form.isbn} onChange={e => setForm({ ...form, isbn: e.target.value })} /></div>
          <div className="form-group"><label>Price</label><input type="number" step="0.01" value={form.bookPrice} onChange={e => setForm({ ...form, bookPrice: parseFloat(e.target.value) })} /></div>
          <div className="form-group"><label>Genre</label><select value={form.genreId} onChange={e => setForm({ ...form, genreId: Number(e.target.value) })}><option value={0}>Select</option>{refData.genres.map(g => <option key={g.id} value={g.id}>{g.text}</option>)}</select></div>
          <div className="form-group"><label>Publisher</label><select value={form.publisherId} onChange={e => setForm({ ...form, publisherId: Number(e.target.value) })}><option value={0}>Select</option>{refData.publishers.map(p => <option key={p.id} value={p.id}>{p.text}</option>)}</select></div>
          <div className="form-group"><label>Book Type</label><select value={form.bookTypeId} onChange={e => setForm({ ...form, bookTypeId: Number(e.target.value) })}><option value={0}>Select</option>{refData.bookTypes.map(b => <option key={b.id} value={b.id}>{b.text}</option>)}</select></div>
          <div className="form-group"><label>Condition</label><select value={form.conditionId} onChange={e => setForm({ ...form, conditionId: Number(e.target.value) })}><option value={0}>Select</option>{refData.conditions.map(c => <option key={c.id} value={c.id}>{c.text}</option>)}</select></div>
          <button className="btn btn-primary" onClick={submit}>Submit Offer</button>
          <button className="btn" onClick={() => setShowForm(false)} style={{ marginLeft: 8 }}>Cancel</button>
        </div>
      )}
      <h2 style={{ margin: '1rem 0' }}>Your Offers</h2>
      {!offers?.length ? <p>No offers yet.</p> : (
        <table>
          <thead><tr><th>Book</th><th>Author</th><th>Price</th><th>Status</th></tr></thead>
          <tbody>
            {offers.map(o => (
              <tr key={o.id}><td>{o.bookName}</td><td>{o.author}</td><td>${o.bookPrice.toFixed(2)}</td><td>{OfferStatusLabels[o.offerStatus]}</td></tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
