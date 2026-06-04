import { useParams } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getOrder, cancelOrder } from '../services/api';
import { OrderStatus } from '../types';

const statusLabels: Record<number, string> = {
  [OrderStatus.Pending]: 'Pending',
  [OrderStatus.Ordered]: 'Ordered',
  [OrderStatus.Shipped]: 'Shipped',
  [OrderStatus.Delivered]: 'Delivered',
  [OrderStatus.Cancelled]: 'Cancelled',
};

export default function OrderDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();
  const { data: order } = useQuery({ queryKey: ['order', id], queryFn: () => getOrder(Number(id)) });

  const handleCancel = async () => {
    await cancelOrder(Number(id));
    toast.success('Order cancelled');
    queryClient.invalidateQueries({ queryKey: ['order', id] });
  };

  if (!order) return <div>Loading...</div>;

  return (
    <div>
      <h1>Order #{order.id}</h1>
      <p><strong>Status:</strong> {statusLabels[order.orderStatus]}</p>
      <p><strong>Delivery Date:</strong> {new Date(order.deliveryDate).toLocaleDateString()}</p>
      <p><strong>Address:</strong> {order.address?.addressLine1}, {order.address?.city}, {order.address?.state} {order.address?.zipCode}</p>

      <h2>Items</h2>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead><tr><th>Book</th><th>Price</th><th>Quantity</th></tr></thead>
        <tbody>
          {order.orderItems?.map(item => (
            <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>{item.book?.name}</td>
              <td>${item.book?.price.toFixed(2)}</td>
              <td>{item.quantity}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <p><strong>Subtotal:</strong> ${order.subTotal.toFixed(2)}</p>
      <p><strong>Tax:</strong> ${order.tax.toFixed(2)}</p>
      <p><strong>Total:</strong> ${order.total.toFixed(2)}</p>

      {order.orderStatus === OrderStatus.Pending && (
        <button onClick={handleCancel} style={{ marginTop: '1rem', color: 'red' }}>Cancel Order</button>
      )}
    </div>
  );
}
