import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { addressService } from '../services';
import { Address } from '../types';
import { useAuth } from '../context/AuthContext';
import toast from 'react-hot-toast';

export default function Addresses() {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const { data: addresses, isLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: () => addressService.list().then(r => r.data),
    enabled: !!user,
  });

  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  const [form, setForm] = useState({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' });

  const resetForm = () => { setForm({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' }); setEditId(null); setShowForm(false); };

  const startEdit = (a: Address) => {
    setForm({ addressLine1: a.addressLine1, addressLine2: a.addressLine2, city: a.city, state: a.state, country: a.country, zipCode: a.zipCode });
    setEditId(a.id);
    setShowForm(true);
  };

  const save = async () => {
    try {
      if (editId) {
        await addressService.update(editId, form);
        toast.success('Address updated');
      } else {
        await addressService.create(form);
        toast.success('Address created');
      }
      queryClient.invalidateQueries({ queryKey: ['addresses'] });
      resetForm();
    } catch { toast.error('Failed to save'); }
  };

  const remove = async (id: number) => {
    try {
      await addressService.delete(id);
      queryClient.invalidateQueries({ queryKey: ['addresses'] });
      toast.success('Address deleted');
    } catch { toast.error('Failed to delete'); }
  };

  if (!user) return <p>Please log in.</p>;
  if (isLoading) return <p>Loading...</p>;

  return (
    <div>
      <h1 className="page-title">My Addresses</h1>
      <button className="btn btn-primary" onClick={() => setShowForm(true)} style={{ marginBottom: '1rem' }}>Add Address</button>
      {showForm && (
        <div className="card">
          <div className="form-group"><label>Address Line 1</label><input value={form.addressLine1} onChange={e => setForm({ ...form, addressLine1: e.target.value })} /></div>
          <div className="form-group"><label>Address Line 2</label><input value={form.addressLine2} onChange={e => setForm({ ...form, addressLine2: e.target.value })} /></div>
          <div className="form-group"><label>City</label><input value={form.city} onChange={e => setForm({ ...form, city: e.target.value })} /></div>
          <div className="form-group"><label>State</label><input value={form.state} onChange={e => setForm({ ...form, state: e.target.value })} /></div>
          <div className="form-group"><label>Country</label><input value={form.country} onChange={e => setForm({ ...form, country: e.target.value })} /></div>
          <div className="form-group"><label>Zip Code</label><input value={form.zipCode} onChange={e => setForm({ ...form, zipCode: e.target.value })} /></div>
          <button className="btn btn-primary" onClick={save}>{editId ? 'Update' : 'Create'}</button>
          <button className="btn" onClick={resetForm} style={{ marginLeft: 8 }}>Cancel</button>
        </div>
      )}
      {addresses?.map(a => (
        <div key={a.id} className="card">
          <p>{a.addressLine1}{a.addressLine2 ? `, ${a.addressLine2}` : ''}</p>
          <p>{a.city}, {a.state} {a.zipCode}, {a.country}</p>
          <button className="btn" onClick={() => startEdit(a)}>Edit</button>
          <button className="btn btn-danger" onClick={() => remove(a.id)} style={{ marginLeft: 4 }}>Delete</button>
        </div>
      ))}
    </div>
  );
}
