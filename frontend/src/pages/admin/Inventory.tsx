import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getAdminInventory, createBook, updateBook } from '../../services/api'
import { getReferenceData } from '../../services/api'
import toast from 'react-hot-toast'
import type { Book } from '../../types'

function Inventory() {
  const queryClient = useQueryClient()
  const [pageIndex, setPageIndex] = useState(1)
  const [editing, setEditing] = useState<Book | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ name: '', author: '', isbn: '', year: '', bookTypeId: 0, conditionId: 0, genreId: 0, publisherId: 0, price: 0, quantity: 0, summary: '' })

  const { data } = useQuery({
    queryKey: ['adminInventory', pageIndex],
    queryFn: async () => { const res = await getAdminInventory({ pageIndex, pageSize: 10 }); return res.data; },
  })

  const { data: refData } = useQuery({
    queryKey: ['referenceData'],
    queryFn: async () => { const res = await getReferenceData(); return res.data; },
  })

  const createMutation = useMutation({
    mutationFn: () => createBook({ ...form, year: form.year ? Number(form.year) : null, id: undefined }),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminInventory'] }); setShowForm(false); toast.success('Book created'); },
  })

  const updateMutation = useMutation({
    mutationFn: () => updateBook({ ...form, id: editing!.id, year: form.year ? Number(form.year) : null }),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminInventory'] }); setShowForm(false); setEditing(null); toast.success('Book updated'); },
  })

  function startEdit(book: Book) {
    setEditing(book)
    setForm({ name: book.name, author: book.author, isbn: book.isbn, year: book.year?.toString() || '', bookTypeId: book.bookTypeId, conditionId: book.conditionId, genreId: book.genreId, publisherId: book.publisherId, price: book.price, quantity: book.quantity, summary: book.summary })
    setShowForm(true)
  }

  return (
    <div>
      <h1>Inventory Management</h1>
      <button onClick={() => { setEditing(null); setForm({ name: '', author: '', isbn: '', year: '', bookTypeId: 0, conditionId: 0, genreId: 0, publisherId: 0, price: 0, quantity: 0, summary: '' }); setShowForm(true); }}>Add Book</button>

      {showForm && refData && (
        <div style={{ margin: '1rem 0', padding: '1rem', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h3>{editing ? 'Edit Book' : 'Add Book'}</h3>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '0.5rem' }}>
            <input placeholder="Name" value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} />
            <input placeholder="Author" value={form.author} onChange={e => setForm({ ...form, author: e.target.value })} />
            <input placeholder="ISBN" value={form.isbn} onChange={e => setForm({ ...form, isbn: e.target.value })} />
            <input placeholder="Year" type="number" value={form.year} onChange={e => setForm({ ...form, year: e.target.value })} />
            <input placeholder="Price" type="number" step="0.01" value={form.price || ''} onChange={e => setForm({ ...form, price: parseFloat(e.target.value) || 0 })} />
            <input placeholder="Quantity" type="number" value={form.quantity || ''} onChange={e => setForm({ ...form, quantity: parseInt(e.target.value) || 0 })} />
            <select value={form.bookTypeId} onChange={e => setForm({ ...form, bookTypeId: Number(e.target.value) })}><option value={0}>Book Type</option>{refData.bookTypes.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}</select>
            <select value={form.conditionId} onChange={e => setForm({ ...form, conditionId: Number(e.target.value) })}><option value={0}>Condition</option>{refData.conditions.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}</select>
            <select value={form.genreId} onChange={e => setForm({ ...form, genreId: Number(e.target.value) })}><option value={0}>Genre</option>{refData.genres.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}</select>
            <select value={form.publisherId} onChange={e => setForm({ ...form, publisherId: Number(e.target.value) })}><option value={0}>Publisher</option>{refData.publishers.map(t => <option key={t.id} value={t.id}>{t.text}</option>)}</select>
          </div>
          <textarea placeholder="Summary" value={form.summary} onChange={e => setForm({ ...form, summary: e.target.value })} style={{ width: '100%', marginTop: '0.5rem' }} />
          <div style={{ marginTop: '0.5rem' }}>
            <button onClick={() => editing ? updateMutation.mutate() : createMutation.mutate()}>{editing ? 'Update' : 'Create'}</button>
            <button onClick={() => setShowForm(false)} style={{ marginLeft: '0.5rem' }}>Cancel</button>
          </div>
        </div>
      )}

      {data && (
        <>
          <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: '1rem' }}>
            <thead><tr><th style={{ textAlign: 'left' }}>Name</th><th style={{ textAlign: 'left' }}>Author</th><th>Price</th><th>Stock</th><th>Actions</th></tr></thead>
            <tbody>
              {data.items.map(book => (
                <tr key={book.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>{book.name}</td>
                  <td style={{ padding: '0.5rem' }}>{book.author}</td>
                  <td style={{ padding: '0.5rem' }}>${book.price.toFixed(2)}</td>
                  <td style={{ padding: '0.5rem', color: book.quantity <= 5 ? 'red' : 'inherit' }}>{book.quantity}</td>
                  <td style={{ padding: '0.5rem' }}><button onClick={() => startEdit(book)}>Edit</button></td>
                </tr>
              ))}
            </tbody>
          </table>
          {data.totalPages > 1 && (
            <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
              <button onClick={() => setPageIndex(p => Math.max(1, p - 1))} disabled={pageIndex === 1}>Previous</button>
              <span>Page {pageIndex} of {data.totalPages}</span>
              <button onClick={() => setPageIndex(p => p + 1)} disabled={pageIndex === data.totalPages}>Next</button>
            </div>
          )}
        </>
      )}
    </div>
  )
}

export default Inventory
