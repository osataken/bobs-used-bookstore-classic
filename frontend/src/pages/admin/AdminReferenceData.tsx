import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getAdminReferenceData, createReferenceData, updateReferenceData } from '../../services/api';
import { ReferenceDataType } from '../../types';
import type { ReferenceDataItem } from '../../types';

const typeLabels: Record<number, string> = {
  [ReferenceDataType.Publisher]: 'Publisher',
  [ReferenceDataType.Condition]: 'Condition',
  [ReferenceDataType.BookType]: 'Book Type',
  [ReferenceDataType.Genre]: 'Genre',
};

export default function AdminReferenceData() {
  const [pageIndex, setPageIndex] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [editItem, setEditItem] = useState<ReferenceDataItem | null>(null);
  const [form, setForm] = useState({ dataType: 0, text: '' });
  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: ['admin-reference-data', pageIndex],
    queryFn: () => getAdminReferenceData({ pageIndex, pageSize: 10 }),
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (editItem) {
      await updateReferenceData(editItem.id, form);
      toast.success('Reference data updated');
    } else {
      await createReferenceData(form);
      toast.success('Reference data created');
    }
    setShowForm(false);
    setEditItem(null);
    setForm({ dataType: 0, text: '' });
    queryClient.invalidateQueries({ queryKey: ['admin-reference-data'] });
  };

  const handleEdit = (item: ReferenceDataItem) => {
    setEditItem(item);
    setForm({ dataType: item.dataType, text: item.text });
    setShowForm(true);
  };

  return (
    <div>
      <h1>Reference Data</h1>
      <button onClick={() => { setShowForm(true); setEditItem(null); setForm({ dataType: 0, text: '' }); }} style={{ marginBottom: '1rem' }}>
        Add New
      </button>

      {showForm && (
        <form onSubmit={handleSubmit} style={{ display: 'flex', gap: '0.5rem', marginBottom: '1rem', padding: '1rem', border: '1px solid #ddd' }}>
          <select value={form.dataType} onChange={e => setForm({ ...form, dataType: Number(e.target.value) })}>
            {Object.entries(typeLabels).map(([val, label]) => (
              <option key={val} value={val}>{label}</option>
            ))}
          </select>
          <input required placeholder="Text" value={form.text} onChange={e => setForm({ ...form, text: e.target.value })} />
          <button type="submit">{editItem ? 'Update' : 'Create'}</button>
          <button type="button" onClick={() => setShowForm(false)}>Cancel</button>
        </form>
      )}

      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead><tr><th>ID</th><th>Type</th><th>Text</th><th>Actions</th></tr></thead>
        <tbody>
          {(data?.items as ReferenceDataItem[])?.map(item => (
            <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>{item.id}</td>
              <td>{typeLabels[item.dataType]}</td>
              <td>{item.text}</td>
              <td><button onClick={() => handleEdit(item)}>Edit</button></td>
            </tr>
          ))}
        </tbody>
      </table>
      {data && data.totalPages > 1 && (
        <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
          <button disabled={pageIndex <= 1} onClick={() => setPageIndex(p => p - 1)}>Previous</button>
          <span>Page {pageIndex} of {data.totalPages}</span>
          <button disabled={pageIndex >= data.totalPages} onClick={() => setPageIndex(p => p + 1)}>Next</button>
        </div>
      )}
    </div>
  );
}
