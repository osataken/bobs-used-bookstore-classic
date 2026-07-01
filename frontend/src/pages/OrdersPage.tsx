import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import api from '../services/api';
import { Order, OrderStatusLabels } from '../types';

function OrdersPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['orders'],
    queryFn: async () => {
      const res = await api.get('/orders');
      return res.data.orders as Order[];
    },
  });

  const cancelOrder = async (orderId: number) => {
    try {
      await api.post('/orders/cancel', { orderId });
      toast.success('Order cancelled');
    } catch {
      toast.error('Failed to cancel order');
    }
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>My Orders</h1>
      {!data || data.length === 0 ? (
        <p>No orders yet.</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Order #</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Status</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Total</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Delivery Date</th>
              <th style={{ padding: '0.5rem' }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {data.map((order) => (
              <tr key={order.id} style={{ borderBottom: '1px solid #ddd' }}>
                <td style={{ padding: '0.5rem' }}><Link to={`/orders/${order.id}`}>#{order.id}</Link></td>
                <td style={{ padding: '0.5rem' }}>{OrderStatusLabels[order.orderStatus]}</td>
                <td style={{ padding: '0.5rem' }}>${order.total.toFixed(2)}</td>
                <td style={{ padding: '0.5rem' }}>{order.deliveryDate}</td>
                <td style={{ padding: '0.5rem' }}>
                  {order.orderStatus === 0 && (
                    <button onClick={() => cancelOrder(order.id)}>Cancel</button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default OrdersPage;
