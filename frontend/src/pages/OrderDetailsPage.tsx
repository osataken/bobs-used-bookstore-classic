import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import api from '../services/api';
import { Order, OrderStatusLabels } from '../types';

function OrderDetailsPage() {
  const { id } = useParams<{ id: string }>();

  const { data: order, isLoading } = useQuery({
    queryKey: ['order', id],
    queryFn: async () => {
      const res = await api.get(`/orders/${id}`);
      return res.data as Order;
    },
  });

  if (isLoading) return <div>Loading...</div>;
  if (!order) return <div>Order not found</div>;

  return (
    <div>
      <h1>Order #{order.id}</h1>
      <p><strong>Status:</strong> {OrderStatusLabels[order.orderStatus]}</p>
      <p><strong>Delivery Date:</strong> {new Date(order.deliveryDate).toLocaleDateString()}</p>
      <p><strong>Subtotal:</strong> ${order.subTotal.toFixed(2)}</p>
      <p><strong>Tax:</strong> ${order.tax.toFixed(2)}</p>
      <p><strong>Total:</strong> ${order.total.toFixed(2)}</p>

      {order.address && (
        <div>
          <h2>Delivery Address</h2>
          <p>{order.address.addressLine1}</p>
          {order.address.addressLine2 && <p>{order.address.addressLine2}</p>}
          <p>{order.address.city}, {order.address.state} {order.address.zipCode}</p>
          <p>{order.address.country}</p>
        </div>
      )}

      {order.orderItems && (
        <div>
          <h2>Items</h2>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Book</th>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Price</th>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Qty</th>
              </tr>
            </thead>
            <tbody>
              {order.orderItems.map((item) => (
                <tr key={item.id}>
                  <td style={{ padding: '0.5rem' }}>{item.book.name}</td>
                  <td style={{ padding: '0.5rem' }}>${item.book.price.toFixed(2)}</td>
                  <td style={{ padding: '0.5rem' }}>{item.quantity}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default OrderDetailsPage;
