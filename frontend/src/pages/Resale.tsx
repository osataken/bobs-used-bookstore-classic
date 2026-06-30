import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getOffers, createOffer, getReferenceData } from '../services/api'
import toast from 'react-hot-toast'

function Resale() {
  const queryClient = useQueryClient()
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ bookName: '', author: '', isbn: '', bookTypeId: 0, conditionId: 0, genreId: 0, publisherId: 0, bookPrice: 0, summary: '' })

  const { data: offers } = useQuery({
    queryKey: ['myOffers'],
    queryFn: async () => { const res = await getOffers(); return res.data.offers; },
  })

  const { data: refData } = useQuery({
    queryKey: ['referenceData'],
    queryFn: async () => { const res = await getReferenceData(); return res.data; },
  })

  const createMutation = useMutation({
    mutationFn: () => createOffer(form),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['myOffers'] }); setShowForm(false); toast.success('Offer submitted'); },
  })

  return (
    <div>
      <h1>Sell Your Books</h1>
      <button onClick={() => setShowForm(true)}>Submit New Offer</button>
      {showForm && refData && (
        <div style={{ margin: '1rem 0', padding: '1rem', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h3>New Offer</h3>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem' }}>
            <input placeholder="Book Name" value={form.bookName} onChange={e => setForm({ ...form, bookName: e.target.value })} />
            <input placeholder="Author" value={form.author} onChange={e => setForm({ ...form, author: e.target.value })} />
            <input placeholder="ISBN" value={form.isbn} onChange={e => setForm({ ...form, isbn: e.target.value })} />
            <input placeholder="Price" type="number" step="0.01" value={form.bookPrice || ''} onChange={e => setForm({ ...form, bookPrice: parseFloat(e.target.value) || 0 })} />
            <select value={form.bookTypeId} onChange={e => setForm({ ...form, bookTypeId: Number(e.target.value) })}>
              <option value={0}>Select Book Type</option>
              {refData.bookTypes.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}
            </select>
            <select value={form.conditionId} onChange={e => setForm({ ...form, conditionId: Number(e.target.value) })}>
              <option value={0}>Select Condition</option>
              {refData.conditions.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}
            </select>
            <select value={form.genreId} onChange={e => setForm({ ...form, genreId: Number(e.target.value) })}>
              <option value={0}>Select Genre</option>
              {refData.genres.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}
            </select>
            <select value={form.publisherId} onChange={e => setForm({ ...form, publisherId: Number(e.target.value) })}>
              <option value={0}>Select Publisher</option>
              {refData.publishers.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}
            </select>
          </div>
          <textarea placeholder="Summary" value={form.summary} onChange={e => setForm({ ...form, summary: e.target.value })} style={{ width: '100%', marginTop: '0.5rem' }} />
          <div style={{ marginTop: '0.5rem' }}>
            <button onClick={() => createMutation.mutate()}>Submit Offer</button>
            <button onClick={() => setShowForm(false)} style={{ marginLeft: '0.5rem' }}>Cancel</button>
          </div>
        </div>
      )}
      <h2>My Offers</h2>
      {(!offers || offers.length === 0) ? <p>No offers submitted.</p> : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead><tr><th style={{ textAlign: 'left' }}>Book</th><th style={{ textAlign: 'left' }}>Author</th><th style={{ textAlign: 'left' }}>Price</th><th style={{ textAlign: 'left' }}>Status</th></tr></thead>
          <tbody>
            {offers.map(offer => (
              <tr key={offer.id} style={{ borderBottom: '1px solid #ddd' }}>
                <td style={{ padding: '0.5rem' }}>{offer.bookName}</td>
                <td style={{ padding: '0.5rem' }}>{offer.author}</td>
                <td style={{ padding: '0.5rem' }}>${offer.bookPrice.toFixed(2)}</td>
                <td style={{ padding: '0.5rem' }}>{offer.offerStatus}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default Resale
