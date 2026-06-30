import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getAdminReferenceData, createReferenceData, updateReferenceData } from '../../services/api'
import toast from 'react-hot-toast'
import type { ReferenceDataItem } from '../../types'

function AdminReferenceData() {
  const queryClient = useQueryClient()
  const [pageIndex, setPageIndex] = useState(1)
  const [editing, setEditing] = useState<ReferenceDataItem | null>(null)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ dataType: 0, text: '' })

  const { data } = useQuery({
    queryKey: ['adminRefData', pageIndex],
    queryFn: async () => { const res = await getAdminReferenceData({ pageIndex, pageSize: 10 }); return res.data; },
  })

  const createMutation = useMutation({
    mutationFn: () => createReferenceData(form),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminRefData'] }); setShowForm(false); toast.success('Created'); },
  })

  const updateMutation = useMutation({
    mutationFn: () => updateReferenceData({ id: editing!.id, ...form }),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminRefData'] }); setShowForm(false); setEditing(null); toast.success('Updated'); },
  })

  const typeLabels = ['Publisher', 'Condition', 'BookType', 'Genre']

  return (
    <div>
      <h1>Reference Data Management</h1>
      <button onClick={() => { setEditing(null); setForm({ dataType: 0, text: '' }); setShowForm(true); }}>Add New</button>

      {showForm && (
        <div style={{ margin: '1rem 0', padding: '1rem', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h3>{editing ? 'Edit' : 'New'} Reference Data</h3>
          <select value={form.dataType} onChange={e => setForm({ ...form, dataType: Number(e.target.value) })}>
            {typeLabels.map((label, i) => <option key={i} value={i}>{label}</option>)}
          </select>
          <input placeholder="Text" value={form.text} onChange={e => setForm({ ...form, text: e.target.value })} style={{ marginLeft: '0.5rem' }} />
          <button onClick={() => editing ? updateMutation.mutate() : createMutation.mutate()} style={{ marginLeft: '0.5rem' }}>{editing ? 'Update' : 'Create'}</button>
          <button onClick={() => setShowForm(false)} style={{ marginLeft: '0.5rem' }}>Cancel</button>
        </div>
      )}

      {data && (
        <>
          <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: '1rem' }}>
            <thead><tr><th style={{ textAlign: 'left' }}>ID</th><th style={{ textAlign: 'left' }}>Type</th><th style={{ textAlign: 'left' }}>Text</th><th>Actions</th></tr></thead>
            <tbody>
              {data.items.map(item => (
                <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>{item.id}</td>
                  <td style={{ padding: '0.5rem' }}>{typeLabels[item.dataType]}</td>
                  <td style={{ padding: '0.5rem' }}>{item.text}</td>
                  <td style={{ padding: '0.5rem' }}><button onClick={() => { setEditing(item); setForm({ dataType: item.dataType, text: item.text }); setShowForm(true); }}>Edit</button></td>
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

export default AdminReferenceData
