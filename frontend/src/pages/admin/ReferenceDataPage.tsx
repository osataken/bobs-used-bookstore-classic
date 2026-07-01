import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import api from '../../services/api';
import { ReferenceDataItem, PaginatedResponse, ReferenceDataTypeLabels } from '../../types';
import Pagination from '../../components/Pagination';

function ReferenceDataPage() {
  const queryClient = useQueryClient();
  const [pageIndex, setPageIndex] = useState(1);
  const [dataTypeFilter, setDataTypeFilter] = useState('');
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ dataType: '', text: '' });

  const { data, isLoading } = useQuery({
    queryKey: ['adminRefData', pageIndex, dataTypeFilter],
    queryFn: async () => {
      const params: Record<string, string | number> = { pageIndex, pageSize: 10 };
      if (dataTypeFilter) params.dataType = dataTypeFilter;
      const res = await api.get('/admin/referenceData', { params });
      return res.data as PaginatedResponse<ReferenceDataItem>;
    },
  });

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/admin/referenceData', { dataType: Number(form.dataType), text: form.text });
      toast.success('Reference data created');
      setShowForm(false);
      setForm({ dataType: '', text: '' });
      queryClient.invalidateQueries({ queryKey: ['adminRefData'] });
    } catch {
      toast.error('Failed to create');
    }
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Reference Data Management</h1>
      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
        <select value={dataTypeFilter} onChange={(e) => { setDataTypeFilter(e.target.value); setPageIndex(1); }}>
          <option value="">All Types</option>
          {Object.entries(ReferenceDataTypeLabels).map(([val, label]) => (
            <option key={val} value={val}>{label}</option>
          ))}
        </select>
        <button onClick={() => setShowForm(!showForm)}>{showForm ? 'Cancel' : 'Add New'}</button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} style={{ margin: '1rem 0', display: 'flex', gap: '0.5rem' }}>
          <select required value={form.dataType} onChange={(e) => setForm({ ...form, dataType: e.target.value })}>
            <option value="">Select Type</option>
            {Object.entries(ReferenceDataTypeLabels).map(([val, label]) => (
              <option key={val} value={val}>{label}</option>
            ))}
          </select>
          <input placeholder="Text" required value={form.text} onChange={(e) => setForm({ ...form, text: e.target.value })} />
          <button type="submit">Create</button>
        </form>
      )}

      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>ID</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Type</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Text</th>
          </tr>
        </thead>
        <tbody>
          {data?.items.map((item) => (
            <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>{item.id}</td>
              <td style={{ padding: '0.5rem' }}>{ReferenceDataTypeLabels[item.dataType]}</td>
              <td style={{ padding: '0.5rem' }}>{item.text}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
    </div>
  );
}

export default ReferenceDataPage;
