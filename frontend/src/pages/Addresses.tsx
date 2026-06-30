import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getAddresses, createAddress, updateAddress, deleteAddress } from '../services/api'
import toast from 'react-hot-toast'
import type { Address } from '../types'

function Addresses() {
  const queryClient = useQueryClient()
  const [editing, setEditing] = useState<Address | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' })

  const { data, isLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: async () => { const res = await getAddresses(); return res.data.addresses; },
  })

  const createMutation = useMutation({
    mutationFn: () => createAddress(form),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['addresses'] }); setShowForm(false); toast.success('Address created'); resetForm(); },
  })

  const updateMutation = useMutation({
    mutationFn: () => updateAddress({ id: editing!.id, ...form }),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['addresses'] }); setEditing(null); toast.success('Address updated'); resetForm(); },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteAddress(id),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['addresses'] }); toast.success('Address deleted'); },
  })

  function resetForm() { setForm({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' }); }
  function startEdit(addr: Address) { setEditing(addr); setForm({ addressLine1: addr.addressLine1, addressLine2: addr.addressLine2, city: addr.city, state: addr.state, country: addr.country, zipCode: addr.zipCode }); setShowForm(true); }

  if (isLoading) return <p>Loading...</p>

  return (
    <div>
      <h1>My Addresses</h1>
      <button onClick={() => { setEditing(null); resetForm(); setShowForm(true); }}>Add Address</button>
      {showForm && (
        <div style={{ margin: '1rem 0', padding: '1rem', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h3>{editing ? 'Edit Address' : 'New Address'}</h3>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem' }}>
            <input placeholder="Address Line 1" value={form.addressLine1} onChange={e => setForm({ ...form, addressLine1: e.target.value })} />
            <input placeholder="Address Line 2" value={form.addressLine2} onChange={e => setForm({ ...form, addressLine2: e.target.value })} />
            <input placeholder="City" value={form.city} onChange={e => setForm({ ...form, city: e.target.value })} />
            <input placeholder="State" value={form.state} onChange={e => setForm({ ...form, state: e.target.value })} />
            <input placeholder="Country" value={form.country} onChange={e => setForm({ ...form, country: e.target.value })} />
            <input placeholder="Zip Code" value={form.zipCode} onChange={e => setForm({ ...form, zipCode: e.target.value })} />
          </div>
          <div style={{ marginTop: '0.5rem', display: 'flex', gap: '0.5rem' }}>
            <button onClick={() => editing ? updateMutation.mutate() : createMutation.mutate()}>
              {editing ? 'Update' : 'Create'}
            </button>
            <button onClick={() => { setShowForm(false); setEditing(null); }}>Cancel</button>
          </div>
        </div>
      )}
      {data && data.length > 0 && (
        <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: '1rem' }}>
          <thead><tr><th style={{ textAlign: 'left' }}>Address</th><th>Actions</th></tr></thead>
          <tbody>
            {data.map(addr => (
              <tr key={addr.id} style={{ borderBottom: '1px solid #ddd' }}>
                <td style={{ padding: '0.5rem' }}>{addr.addressLine1}, {addr.city}, {addr.state} {addr.zipCode}, {addr.country}</td>
                <td style={{ padding: '0.5rem', display: 'flex', gap: '0.5rem' }}>
                  <button onClick={() => startEdit(addr)}>Edit</button>
                  <button onClick={() => deleteMutation.mutate(addr.id)}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default Addresses
