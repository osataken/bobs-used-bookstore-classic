import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import api from '../services/api';
import { Address } from '../types';

function AddressPage() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [editId, setEditId] = useState<number | null>(null);
  const [form, setForm] = useState({
    addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '',
  });

  const { data: addresses } = useQuery({
    queryKey: ['addresses'],
    queryFn: async () => {
      const res = await api.get('/addresses');
      return res.data.addresses as Address[];
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editId) {
        await api.put('/addresses', { id: editId, ...form });
        toast.success('Address updated');
      } else {
        await api.post('/addresses', form);
        toast.success('Address created');
      }
      setShowForm(false);
      setEditId(null);
      setForm({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' });
      queryClient.invalidateQueries({ queryKey: ['addresses'] });
    } catch {
      toast.error('Failed to save address');
    }
  };

  const handleEdit = (addr: Address) => {
    setForm({
      addressLine1: addr.addressLine1, addressLine2: addr.addressLine2,
      city: addr.city, state: addr.state, country: addr.country, zipCode: addr.zipCode,
    });
    setEditId(addr.id);
    setShowForm(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await api.post('/addresses/delete', { id });
      toast.success('Address deleted');
      queryClient.invalidateQueries({ queryKey: ['addresses'] });
    } catch {
      toast.error('Failed to delete address');
    }
  };

  return (
    <div>
      <h1>My Addresses</h1>
      <button onClick={() => { setShowForm(!showForm); setEditId(null); setForm({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' }); }}>
        {showForm ? 'Cancel' : 'Add New Address'}
      </button>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ margin: '1rem 0', display: 'grid', gap: '0.5rem', maxWidth: '400px' }}>
          <input placeholder="Address Line 1" required value={form.addressLine1} onChange={(e) => setForm({ ...form, addressLine1: e.target.value })} />
          <input placeholder="Address Line 2" value={form.addressLine2} onChange={(e) => setForm({ ...form, addressLine2: e.target.value })} />
          <input placeholder="City" required value={form.city} onChange={(e) => setForm({ ...form, city: e.target.value })} />
          <input placeholder="State" required value={form.state} onChange={(e) => setForm({ ...form, state: e.target.value })} />
          <input placeholder="Country" required value={form.country} onChange={(e) => setForm({ ...form, country: e.target.value })} />
          <input placeholder="Zip Code" required value={form.zipCode} onChange={(e) => setForm({ ...form, zipCode: e.target.value })} />
          <button type="submit">{editId ? 'Update' : 'Create'}</button>
        </form>
      )}

      {!addresses || addresses.length === 0 ? (
        <p>No addresses saved.</p>
      ) : (
        <div style={{ marginTop: '1rem' }}>
          {addresses.map((addr) => (
            <div key={addr.id} style={{ border: '1px solid #ddd', padding: '1rem', marginBottom: '0.5rem', borderRadius: '8px' }}>
              <p>{addr.addressLine1}</p>
              {addr.addressLine2 && <p>{addr.addressLine2}</p>}
              <p>{addr.city}, {addr.state} {addr.zipCode}</p>
              <p>{addr.country}</p>
              <div style={{ display: 'flex', gap: '0.5rem' }}>
                <button onClick={() => handleEdit(addr)}>Edit</button>
                <button onClick={() => handleDelete(addr.id)}>Delete</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default AddressPage;
