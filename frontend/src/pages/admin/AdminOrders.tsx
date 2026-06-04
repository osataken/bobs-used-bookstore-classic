import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { getAdminOrders } from '../../services/api';
import { OrderStatus } from '../../types';
import type { Order } from '../../types';

const statusLabels: Record<number, string> = {
  [OrderStatus.Pending]: 'Pending',
  [OrderStatus.Ordered]: 'Ordered',
  [OrderStatus.Shipped]: 'Shipped',
  [OrderStatus.Delivered]: 'Delivered',
  [OrderStatus.Cancelled]: 'Cancelled',
};

export default function AdminOrders() {
  const [pageIndex, setPageIndex] = useState(1);

  const { data } = useQuery({
    queryKey: ['admin-orders', pageIndex],
    queryFn: () => getAdminOrders({ pageIndex, pageSize: 10 }),
  });

  return (
    <div>
      <h1>Orders Management</h1>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead><tr><th>Order #</th><th>Customer</th><th>Status</th><th>Date</th><th>Actions</th></tr></thead>
        <tbody>
          {(data?.items as Order[])?.map(order => (
            <tr key={order.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>{order.id}</td>
              <td>{order.customer?.firstName} {order.customer?.lastName}</td>
              <td>{statusLabels[order.orderStatus]}</td>
              <td>{new Date(order.createdOn).toLocaleDateString()}</td>
              <td><Link to={`/admin/orders/${order.id}`}>Details</Link></td>
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
