import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getAdminReferenceData, createReferenceData, updateReferenceData } from '../../services/api'
import Pagination from '../../components/Pagination'
import toast from 'react-hot-toast'
import type { ReferenceData } from '../../types'

export default function ReferenceDataPage() {
  const queryClient = useQueryClient()
  const [pageIndex, setPageIndex] = useState(1)
  const [dataTypeFilter, setDataTypeFilter] = useState<number | undefined>(undefined)
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<ReferenceData | null>(null)
  const [form, setForm] = useState({ dataType: 0, text: '' })

  const { data, isLoading } = useQuery({
    queryKey: ['adminRefData', pageIndex, dataTypeFilter],
    queryFn: () => getAdminReferenceData({ pageIndex, pageSize: 10, dataType: dataTypeFilter }).then(r => r.data),
  })

  const handleEdit = (item: ReferenceData) => {
    setEditing(item)
    setForm({ dataType: item.dataType, text: item.text })
    setShowForm(true)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      if (editing) {
        await updateReferenceData(editing.id, form)
        toast.success('Updated')
      } else {
        await createReferenceData(form)
        toast.success('Created')
      }
      queryClient.invalidateQueries({ queryKey: ['adminRefData'] })
      setShowForm(false)
      setEditing(null)
    } catch {
      toast.error('Failed to save')
    }
  }

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      <h1>Reference Data</h1>
      <div style={{ display: 'flex', gap: '10px', marginBottom: '15px' }}>
        <select value={dataTypeFilter ?? ''} onChange={e => { setDataTypeFilter(e.target.value ? Number(e.target.value) : undefined); setPageIndex(1) }}>
          <option value="">All Types</option>
          <option value={0}>Publisher</option>
          <option value={1}>Condition</option>
          <option value={2}>Book Type</option>
          <option value={3}>Genre</option>
        </select>
        <button onClick={() => { setShowForm(true); setEditing(null); setForm({ dataType: 0, text: '' }) }}>Add New</button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ border: '1px solid #ddd', padding: '15px', marginBottom: '20px' }}>
          <select value={form.dataType} onChange={e => setForm({ ...form, dataType: Number(e.target.value) })}>
            <option value={0}>Publisher</option>
            <option value={1}>Condition</option>
            <option value={2}>Book Type</option>
            <option value={3}>Genre</option>
          </select>
          <input value={form.text} onChange={e => setForm({ ...form, text: e.target.value })} placeholder="Text" required style={{ marginLeft: '10px' }} />
          <button type="submit" style={{ marginLeft: '10px' }}>Save</button>
          <button type="button" onClick={() => setShowForm(false)} style={{ marginLeft: '5px' }}>Cancel</button>
        </form>
      )}

      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>ID</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Type</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Text</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {data?.items?.map(item => (
            <tr key={item.id}>
              <td style={{ padding: '8px' }}>{item.id}</td>
              <td style={{ padding: '8px' }}>{item.typeText}</td>
              <td style={{ padding: '8px' }}>{item.text}</td>
              <td style={{ padding: '8px' }}><button onClick={() => handleEdit(item)}>Edit</button></td>
            </tr>
          ))}
        </tbody>
      </table>
      {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
    </div>
  )
}
