import { useParams } from 'react-router-dom';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import { ordersService } from '../services';
import { OrderStatusLabels } from '../types';
import toast from 'react-hot-toast';

export default function OrderDetails() {
  const { id } = useParams<{ id: string }>();
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ['order', id],
    queryFn: () => ordersService.get(Number(id)).then(r => r.data),
  });

  const cancel = async () => {
    try {
      await ordersService.cancel(Number(id));
      queryClient.invalidateQueries({ queryKey: ['order', id] });
      toast.success('Order cancelled');
    } catch { toast.error('Failed to cancel'); }
  };

  if (isLoading) return <p>Loading...</p>;
  if (!data) return <p>Order not found</p>;

  const order = data.order || data;

  return (
    <div>
      <h1 className="page-title">Order #{order.id}</h1>
      <div className="card">
        <p><strong>Status:</strong> {OrderStatusLabels[order.orderStatus]}</p>
        <p><strong>Delivery Date:</strong> {new Date(order.deliveryDate).toLocaleDateString()}</p>
        <p><strong>Subtotal:</strong> ${(data.subTotal || order.subTotal || 0).toFixed(2)}</p>
        <p><strong>Tax:</strong> ${(data.tax || order.tax || 0).toFixed(2)}</p>
        <p><strong>Total:</strong> ${(data.total || order.total || 0).toFixed(2)}</p>
      </div>
      {order.orderItems?.length > 0 && (
        <table>
          <thead><tr><th>Book</th><th>Quantity</th><th>Price</th></tr></thead>
          <tbody>
            {order.orderItems.map((item: { id: number; book?: { name: string; price: number }; quantity: number }) => (
              <tr key={item.id}>
                <td>{item.book?.name}</td>
                <td>{item.quantity}</td>
                <td>${item.book?.price.toFixed(2)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {order.orderStatus === 0 && (
        <button className="btn btn-danger" onClick={cancel} style={{ marginTop: '1rem' }}>Cancel Order</button>
      )}
    </div>
  );
}
