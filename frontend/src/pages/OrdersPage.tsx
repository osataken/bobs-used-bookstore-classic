import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { getMyOrders } from '../services/api';
import { OrderStatus } from '../types';

const statusLabels: Record<number, string> = {
  [OrderStatus.Pending]: 'Pending',
  [OrderStatus.Ordered]: 'Ordered',
  [OrderStatus.Shipped]: 'Shipped',
  [OrderStatus.Delivered]: 'Delivered',
  [OrderStatus.Cancelled]: 'Cancelled',
};

export default function OrdersPage() {
  const { data: orders } = useQuery({ queryKey: ['my-orders'], queryFn: getMyOrders });

  return (
    <div>
      <h1>My Orders</h1>
      {!orders?.length ? (
        <p>No orders yet.</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead><tr><th>Order #</th><th>Date</th><th>Status</th><th>Total</th><th>Details</th></tr></thead>
          <tbody>
            {orders.map(order => (
              <tr key={order.id} style={{ borderBottom: '1px solid #ddd' }}>
                <td style={{ padding: '0.5rem' }}>{order.id}</td>
                <td>{new Date(order.createdOn).toLocaleDateString()}</td>
                <td>{statusLabels[order.orderStatus]}</td>
                <td>${order.total.toFixed(2)}</td>
                <td><Link to={`/orders/${order.id}`}>View</Link></td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
