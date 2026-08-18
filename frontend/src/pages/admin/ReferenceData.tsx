import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { adminService } from '../../services';
import { ReferenceDataItem, ReferenceDataTypeLabels } from '../../types';
import { useAuth } from '../../context/AuthContext';
import toast from 'react-hot-toast';

export default function ReferenceData() {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [pageIndex, setPageIndex] = useState(1);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ dataType: 0, text: '' });

  const { data, isLoading } = useQuery({
    queryKey: ['admin-refdata', pageIndex],
    queryFn: () => adminService.listReferenceData({ pageIndex, pageSize: 24 }).then(r => r.data),
    enabled: !!user?.isAdmin,
  });

  const create = async () => {
    try {
      await adminService.createReferenceData(form);
      queryClient.invalidateQueries({ queryKey: ['admin-refdata'] });
      toast.success('Created');
      setShowForm(false);
    } catch { toast.error('Failed'); }
  };

  if (!user?.isAdmin) return <p>Access denied.</p>;

  return (
    <div>
      <h1 className="page-title">Reference Data</h1>
      <button className="btn btn-primary" onClick={() => setShowForm(true)} style={{ marginBottom: '1rem' }}>Add New</button>
      {showForm && (
        <div className="card">
          <div className="form-group"><label>Type</label>
            <select value={form.dataType} onChange={e => setForm({ ...form, dataType: Number(e.target.value) })}>
              <option value={0}>Publisher</option><option value={1}>Condition</option><option value={2}>BookType</option><option value={3}>Genre</option>
            </select>
          </div>
          <div className="form-group"><label>Text</label><input value={form.text} onChange={e => setForm({ ...form, text: e.target.value })} /></div>
          <button className="btn btn-primary" onClick={create}>Create</button>
          <button className="btn" onClick={() => setShowForm(false)} style={{ marginLeft: 8 }}>Cancel</button>
        </div>
      )}
      {isLoading ? <p>Loading...</p> : (
        <table>
          <thead><tr><th>ID</th><th>Type</th><th>Text</th></tr></thead>
          <tbody>
            {(data?.items as ReferenceDataItem[])?.map(item => (
              <tr key={item.id}><td>{item.id}</td><td>{ReferenceDataTypeLabels[item.dataType]}</td><td>{item.text}</td></tr>
            ))}
          </tbody>
        </table>
      )}
      {data && data.totalPages > 1 && (
        <div className="pagination">
          <button className="btn" disabled={pageIndex <= 1} onClick={() => setPageIndex(p => p - 1)}>Prev</button>
          <span>Page {data.pageIndex} of {data.totalPages}</span>
          <button className="btn" disabled={pageIndex >= data.totalPages} onClick={() => setPageIndex(p => p + 1)}>Next</button>
        </div>
      )}
    </div>
  );
}
