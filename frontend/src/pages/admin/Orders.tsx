import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { adminService } from '../../services';
import { Order, OrderStatusLabels } from '../../types';
import { useAuth } from '../../context/AuthContext';

export default function AdminOrders() {
  const { user } = useAuth();
  const [pageIndex, setPageIndex] = useState(1);

  const { data, isLoading } = useQuery({
    queryKey: ['admin-orders', pageIndex],
    queryFn: () => adminService.listOrders({ pageIndex, pageSize: 10 }).then(r => r.data),
    enabled: !!user?.isAdmin,
  });

  if (!user?.isAdmin) return <p>Access denied.</p>;

  return (
    <div>
      <h1 className="page-title">Manage Orders</h1>
      {isLoading ? <p>Loading...</p> : (
        <table>
          <thead><tr><th>ID</th><th>Customer</th><th>Status</th><th>Delivery Date</th></tr></thead>
          <tbody>
            {(data?.items as Order[])?.map(o => (
              <tr key={o.id}>
                <td>{o.id}</td>
                <td>{o.customer?.username || `Customer #${o.customerId}`}</td>
                <td>{OrderStatusLabels[o.orderStatus]}</td>
                <td>{new Date(o.deliveryDate).toLocaleDateString()}</td>
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
