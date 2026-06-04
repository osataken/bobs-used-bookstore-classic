import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getAdminOffers, approveOffer, rejectOffer, receivedOffer, paidOffer } from '../../services/api';
import { OfferStatus } from '../../types';
import type { Offer } from '../../types';

const statusLabels: Record<number, string> = {
  [OfferStatus.PendingApproval]: 'Pending Approval',
  [OfferStatus.Approved]: 'Approved',
  [OfferStatus.Received]: 'Received',
  [OfferStatus.Paid]: 'Paid',
  [OfferStatus.Rejected]: 'Rejected',
};

export default function AdminOffers() {
  const [pageIndex, setPageIndex] = useState(1);
  const queryClient = useQueryClient();

  const { data } = useQuery({
    queryKey: ['admin-offers', pageIndex],
    queryFn: () => getAdminOffers({ pageIndex, pageSize: 10 }),
  });

  const handleAction = async (id: number, action: string) => {
    switch (action) {
      case 'approve': await approveOffer(id); break;
      case 'reject': await rejectOffer(id); break;
      case 'received': await receivedOffer(id); break;
      case 'paid': await paidOffer(id); break;
    }
    toast.success(`Offer ${action}d`);
    queryClient.invalidateQueries({ queryKey: ['admin-offers'] });
  };

  return (
    <div>
      <h1>Offers Management</h1>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead><tr><th>Book</th><th>Author</th><th>Price</th><th>Status</th><th>Actions</th></tr></thead>
        <tbody>
          {(data?.items as Offer[])?.map(offer => (
            <tr key={offer.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>{offer.bookName}</td>
              <td>{offer.author}</td>
              <td>${offer.bookPrice.toFixed(2)}</td>
              <td>{statusLabels[offer.offerStatus]}</td>
              <td style={{ display: 'flex', gap: '0.25rem' }}>
                {offer.offerStatus === OfferStatus.PendingApproval && (
                  <>
                    <button onClick={() => handleAction(offer.id, 'approve')}>Approve</button>
                    <button onClick={() => handleAction(offer.id, 'reject')}>Reject</button>
                  </>
                )}
                {offer.offerStatus === OfferStatus.Approved && (
                  <button onClick={() => handleAction(offer.id, 'received')}>Received</button>
                )}
                {offer.offerStatus === OfferStatus.Received && (
                  <button onClick={() => handleAction(offer.id, 'paid')}>Paid</button>
                )}
              </td>
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
