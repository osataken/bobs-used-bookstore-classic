import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { adminService } from '../../services';
import { Offer, OfferStatusLabels } from '../../types';
import { useAuth } from '../../context/AuthContext';
import toast from 'react-hot-toast';

export default function Offers() {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [pageIndex, setPageIndex] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ['admin-offers', pageIndex],
    queryFn: () => adminService.listOffers({ pageIndex, pageSize: 10 }).then(r => r.data),
    enabled: !!user?.isAdmin,
  });

  const action = async (id: number, fn: (id: number) => Promise<unknown>, msg: string) => {
    try { await fn(id); queryClient.invalidateQueries({ queryKey: ['admin-offers'] }); toast.success(msg); }
    catch { toast.error('Failed'); }
  };

  if (!user?.isAdmin) return <p>Access denied.</p>;

  return (
    <div>
      <h1 className="page-title">Manage Offers</h1>
      {isLoading ? <p>Loading...</p> : (
        <table>
          <thead><tr><th>ID</th><th>Book</th><th>Author</th><th>Price</th><th>Status</th><th>Actions</th></tr></thead>
          <tbody>
            {(data?.items as Offer[])?.map(o => (
              <tr key={o.id}>
                <td>{o.id}</td><td>{o.bookName}</td><td>{o.author}</td><td>${o.bookPrice.toFixed(2)}</td>
                <td>{OfferStatusLabels[o.offerStatus]}</td>
                <td>
                  {o.offerStatus === 0 && <><button className="btn btn-success" onClick={() => action(o.id, adminService.approveOffer, 'Approved')}>Approve</button> <button className="btn btn-danger" onClick={() => action(o.id, adminService.rejectOffer, 'Rejected')}>Reject</button></>}
                  {o.offerStatus === 1 && <button className="btn btn-primary" onClick={() => action(o.id, adminService.receivedOffer, 'Received')}>Received</button>}
                  {o.offerStatus === 2 && <button className="btn btn-primary" onClick={() => action(o.id, adminService.paidOffer, 'Paid')}>Paid</button>}
                </td>
              </tr>
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
