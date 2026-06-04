import { useState, useEffect } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import toast from 'react-hot-toast';
import { getAddress, createAddress, updateAddress } from '../services/api';

export default function AddressFormPage() {
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const isEdit = !!id;
  const returnUrl = searchParams.get('returnUrl') || '/addresses';

  const [form, setForm] = useState({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' });

  useEffect(() => {
    if (isEdit) {
      getAddress(Number(id)).then(addr => {
        setForm({ addressLine1: addr.addressLine1, addressLine2: addr.addressLine2 || '', city: addr.city, state: addr.state, country: addr.country, zipCode: addr.zipCode });
      });
    }
  }, [id, isEdit]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (isEdit) {
      await updateAddress(Number(id), form);
      toast.success('Address updated');
    } else {
      await createAddress(form);
      toast.success('Address created');
    }
    navigate(returnUrl);
  };

  return (
    <div>
      <h1>{isEdit ? 'Edit Address' : 'New Address'}</h1>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: '1rem', maxWidth: '500px' }}>
        <input required placeholder="Address Line 1" value={form.addressLine1} onChange={e => setForm({ ...form, addressLine1: e.target.value })} />
        <input placeholder="Address Line 2" value={form.addressLine2} onChange={e => setForm({ ...form, addressLine2: e.target.value })} />
        <input required placeholder="City" value={form.city} onChange={e => setForm({ ...form, city: e.target.value })} />
        <input required placeholder="State" value={form.state} onChange={e => setForm({ ...form, state: e.target.value })} />
        <input required placeholder="Country" value={form.country} onChange={e => setForm({ ...form, country: e.target.value })} />
        <input required placeholder="Zip Code" value={form.zipCode} onChange={e => setForm({ ...form, zipCode: e.target.value })} />
        <button type="submit">{isEdit ? 'Update' : 'Create'}</button>
      </form>
    </div>
  );
}
