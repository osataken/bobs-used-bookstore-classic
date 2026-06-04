import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getAdminOrder, updateOrderStatus } from '../../services/api';
import { OrderStatus } from '../../types';

const statusLabels: Record<number, string> = {
  [OrderStatus.Pending]: 'Pending',
  [OrderStatus.Ordered]: 'Ordered',
  [OrderStatus.Shipped]: 'Shipped',
  [OrderStatus.Delivered]: 'Delivered',
  [OrderStatus.Cancelled]: 'Cancelled',
};

export default function AdminOrderDetails() {
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();
  const { data: order } = useQuery({ queryKey: ['admin-order', id], queryFn: () => getAdminOrder(Number(id)) });
  const [selectedStatus, setSelectedStatus] = useState<number>(-1);

  const handleUpdateStatus = async () => {
    if (selectedStatus < 0) return;
    await updateOrderStatus(Number(id), selectedStatus as OrderStatus);
    toast.success('Order status updated');
    queryClient.invalidateQueries({ queryKey: ['admin-order', id] });
  };

  if (!order) return <div>Loading...</div>;

  return (
    <div>
      <h1>Order #{order.id}</h1>
      <p><strong>Customer:</strong> {order.customer?.firstName} {order.customer?.lastName}</p>
      <p><strong>Status:</strong> {statusLabels[order.orderStatus]}</p>
      <p><strong>Delivery Date:</strong> {new Date(order.deliveryDate).toLocaleDateString()}</p>

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

      <h2>Update Status</h2>
      <div style={{ display: 'flex', gap: '0.5rem' }}>
        <select value={selectedStatus} onChange={e => setSelectedStatus(Number(e.target.value))}>
          <option value={-1}>Select status...</option>
          {Object.entries(statusLabels).map(([val, label]) => (
            <option key={val} value={val}>{label}</option>
          ))}
        </select>
        <button onClick={handleUpdateStatus} disabled={selectedStatus < 0}>Update</button>
      </div>
    </div>
  );
}
