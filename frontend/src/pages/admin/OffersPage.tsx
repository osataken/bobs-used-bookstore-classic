import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import api from '../../services/api';
import { Offer, PaginatedResponse, OfferStatusLabels } from '../../types';
import Pagination from '../../components/Pagination';

function OffersPage() {
  const queryClient = useQueryClient();
  const [pageIndex, setPageIndex] = useState(1);
  const [statusFilter, setStatusFilter] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['adminOffers', pageIndex, statusFilter],
    queryFn: async () => {
      const params: Record<string, string | number> = { pageIndex, pageSize: 10 };
      if (statusFilter) params.status = statusFilter;
      const res = await api.get('/admin/offers', { params });
      return res.data as PaginatedResponse<Offer>;
    },
  });

  const updateStatus = async (id: number, status: number) => {
    try {
      await api.post('/admin/offers/status', { id, status });
      toast.success('Offer status updated');
      queryClient.invalidateQueries({ queryKey: ['adminOffers'] });
    } catch {
      toast.error('Failed to update status');
    }
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Offers Management</h1>
      <select value={statusFilter} onChange={(e) => { setStatusFilter(e.target.value); setPageIndex(1); }}>
        <option value="">All Statuses</option>
        {Object.entries(OfferStatusLabels).map(([val, label]) => (
          <option key={val} value={val}>{label}</option>
        ))}
      </select>

      <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: '1rem' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Book</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Author</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Price</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Status</th>
            <th style={{ padding: '0.5rem' }}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {data?.items.map((offer) => (
            <tr key={offer.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>{offer.bookName}</td>
              <td style={{ padding: '0.5rem' }}>{offer.author}</td>
              <td style={{ padding: '0.5rem' }}>${offer.bookPrice.toFixed(2)}</td>
              <td style={{ padding: '0.5rem' }}>{OfferStatusLabels[offer.offerStatus]}</td>
              <td style={{ padding: '0.5rem' }}>
                {offer.offerStatus === 0 && (
                  <div style={{ display: 'flex', gap: '0.25rem' }}>
                    <button onClick={() => updateStatus(offer.id, 1)}>Approve</button>
                    <button onClick={() => updateStatus(offer.id, 4)}>Reject</button>
                  </div>
                )}
                {offer.offerStatus === 1 && (
                  <button onClick={() => updateStatus(offer.id, 2)}>Mark Received</button>
                )}
                {offer.offerStatus === 2 && (
                  <button onClick={() => updateStatus(offer.id, 3)}>Mark Paid</button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
    </div>
  );
}

export default OffersPage;
