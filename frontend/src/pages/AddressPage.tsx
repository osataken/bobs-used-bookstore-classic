import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getAddresses, createAddress, updateAddress, deleteAddress } from '../services/api'
import toast from 'react-hot-toast'
import type { Address } from '../types'

export default function AddressPage() {
  const queryClient = useQueryClient()
  const [editing, setEditing] = useState<Address | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' })

  const { data, isLoading } = useQuery({
    queryKey: ['addresses'],
    queryFn: () => getAddresses().then(r => r.data),
  })

  const resetForm = () => {
    setForm({ addressLine1: '', addressLine2: '', city: '', state: '', country: '', zipCode: '' })
    setEditing(null)
    setShowForm(false)
  }

  const handleEdit = (addr: Address) => {
    setEditing(addr)
    setForm({ addressLine1: addr.addressLine1, addressLine2: addr.addressLine2, city: addr.city, state: addr.state, country: addr.country, zipCode: addr.zipCode })
    setShowForm(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editing) {
        await updateAddress(editing.id, form)
        toast.success('Address updated')
      } else {
        await createAddress(form)
        toast.success('Address created')
      }
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      resetForm()
    } catch {
      toast.error('Failed to save address')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await deleteAddress(id)
      queryClient.invalidateQueries({ queryKey: ['addresses'] })
      toast.success('Address deleted')
    } catch {
      toast.error('Failed to delete address')
    }
  }

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      <h1>My Addresses</h1>
      <button onClick={() => { resetForm(); setShowForm(true) }} style={{ marginBottom: '15px' }}>Add New Address</button>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ border: '1px solid #ddd', padding: '15px', marginBottom: '20px' }}>
          <h3>{editing ? 'Edit Address' : 'New Address'}</h3>
          <div style={{ display: 'grid', gap: '10px' }}>
            <input placeholder="Address Line 1" value={form.addressLine1} onChange={e => setForm({ ...form, addressLine1: e.target.value })} required />
            <input placeholder="Address Line 2" value={form.addressLine2} onChange={e => setForm({ ...form, addressLine2: e.target.value })} />
            <input placeholder="City" value={form.city} onChange={e => setForm({ ...form, city: e.target.value })} required />
            <input placeholder="State" value={form.state} onChange={e => setForm({ ...form, state: e.target.value })} required />
            <input placeholder="Country" value={form.country} onChange={e => setForm({ ...form, country: e.target.value })} required />
            <input placeholder="Zip Code" value={form.zipCode} onChange={e => setForm({ ...form, zipCode: e.target.value })} required />
          </div>
          <div style={{ marginTop: '10px' }}>
            <button type="submit">Save</button>{' '}
            <button type="button" onClick={resetForm}>Cancel</button>
          </div>
        </form>
      )}

      {(!data?.addresses || data.addresses.length === 0) ? (
        <p>No addresses found.</p>
      ) : (
        <div>
          {data.addresses.map(addr => (
            <div key={addr.id} style={{ border: '1px solid #ddd', padding: '10px', marginBottom: '10px' }}>
              <p>{addr.addressLine1}{addr.addressLine2 && `, ${addr.addressLine2}`}</p>
              <p>{addr.city}, {addr.state} {addr.zipCode}, {addr.country}</p>
              <button onClick={() => handleEdit(addr)}>Edit</button>{' '}
              <button onClick={() => handleDelete(addr.id)}>Delete</button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
