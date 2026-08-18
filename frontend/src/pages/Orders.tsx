import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { ordersService } from '../services';
import { useAuth } from '../context/AuthContext';
import { OrderStatusLabels } from '../types';

export default function Orders() {
  const { user } = useAuth();
  const { data: orders, isLoading } = useQuery({
    queryKey: ['orders'],
    queryFn: () => ordersService.list().then(r => r.data),
    enabled: !!user,
  });

  if (!user) return <p>Please log in to view orders.</p>;
  if (isLoading) return <p>Loading...</p>;

  return (
    <div>
      <h1 className="page-title">My Orders</h1>
      {!orders?.length ? <p>No orders yet.</p> : (
        <table>
          <thead><tr><th>Order #</th><th>Date</th><th>Status</th><th>Total</th><th>Actions</th></tr></thead>
          <tbody>
            {orders.map(order => (
              <tr key={order.id}>
                <td>{order.id}</td>
                <td>{new Date(order.deliveryDate).toLocaleDateString()}</td>
                <td>{OrderStatusLabels[order.orderStatus]}</td>
                <td>${(order.total || 0).toFixed(2)}</td>
                <td><Link to={`/orders/${order.id}`}>Details</Link></td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
