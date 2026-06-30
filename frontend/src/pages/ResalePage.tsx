import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getResaleOffers, createResaleOffer, getReferenceData } from '../services/api'
import toast from 'react-hot-toast'

export default function ResalePage() {
  const queryClient = useQueryClient()
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ bookName: '', author: '', isbn: '', bookTypeId: 0, conditionId: 0, genreId: 0, publisherId: 0, bookPrice: 0, summary: '' })

  const { data } = useQuery({ queryKey: ['resale'], queryFn: () => getResaleOffers().then(r => r.data) })
  const { data: refData } = useQuery({ queryKey: ['referenceData'], queryFn: () => getReferenceData().then(r => r.data) })

  const genres = refData?.items?.filter(r => r.dataType === 3) || []
  const conditions = refData?.items?.filter(r => r.dataType === 1) || []
  const publishers = refData?.items?.filter(r => r.dataType === 0) || []
  const bookTypes = refData?.items?.filter(r => r.dataType === 2) || []

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await createResaleOffer(form)
      queryClient.invalidateQueries({ queryKey: ['resale'] })
      toast.success('Offer submitted')
      setShowForm(false)
    } catch {
      toast.error('Failed to submit offer')
    }
  }

  return (
    <div>
      <h1>Resale Offers</h1>
      <button onClick={() => setShowForm(!showForm)} style={{ marginBottom: '15px' }}>Submit New Offer</button>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ border: '1px solid #ddd', padding: '15px', marginBottom: '20px' }}>
          <div style={{ display: 'grid', gap: '10px' }}>
            <input placeholder="Book Name" value={form.bookName} onChange={e => setForm({ ...form, bookName: e.target.value })} required />
            <input placeholder="Author" value={form.author} onChange={e => setForm({ ...form, author: e.target.value })} />
            <input placeholder="ISBN" value={form.isbn} onChange={e => setForm({ ...form, isbn: e.target.value })} />
            <input placeholder="Price" type="number" step="0.01" value={form.bookPrice || ''} onChange={e => setForm({ ...form, bookPrice: parseFloat(e.target.value) || 0 })} required />
            <select value={form.genreId} onChange={e => setForm({ ...form, genreId: Number(e.target.value) })} required>
              <option value={0}>Select Genre</option>
              {genres.map(g => <option key={g.id} value={g.id}>{g.text}</option>)}
            </select>
            <select value={form.conditionId} onChange={e => setForm({ ...form, conditionId: Number(e.target.value) })} required>
              <option value={0}>Select Condition</option>
              {conditions.map(c => <option key={c.id} value={c.id}>{c.text}</option>)}
            </select>
            <select value={form.bookTypeId} onChange={e => setForm({ ...form, bookTypeId: Number(e.target.value) })} required>
              <option value={0}>Select Book Type</option>
              {bookTypes.map(bt => <option key={bt.id} value={bt.id}>{bt.text}</option>)}
            </select>
            <select value={form.publisherId} onChange={e => setForm({ ...form, publisherId: Number(e.target.value) })} required>
              <option value={0}>Select Publisher</option>
              {publishers.map(p => <option key={p.id} value={p.id}>{p.text}</option>)}
            </select>
            <textarea placeholder="Summary" value={form.summary} onChange={e => setForm({ ...form, summary: e.target.value })} />
          </div>
          <button type="submit" style={{ marginTop: '10px' }}>Submit Offer</button>
        </form>
      )}

      {(!data?.offers || data.offers.length === 0) ? <p>No offers submitted.</p> : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Book</th>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Price</th>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Status</th>
            </tr>
          </thead>
          <tbody>
            {data.offers.map(offer => (
              <tr key={offer.id}>
                <td style={{ padding: '8px' }}>{offer.bookName}</td>
                <td style={{ padding: '8px' }}>${offer.bookPrice.toFixed(2)}</td>
                <td style={{ padding: '8px' }}>{offer.statusText}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
