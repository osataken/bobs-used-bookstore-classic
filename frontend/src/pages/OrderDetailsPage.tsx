import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getOrderDetails } from '../services/api'

export default function OrderDetailsPage() {
  const { id } = useParams<{ id: string }>()
  const { data: order, isLoading } = useQuery({
    queryKey: ['order', id],
    queryFn: () => getOrderDetails(Number(id)).then(r => r.data),
    enabled: !!id,
  })

  if (isLoading) return <div>Loading...</div>
  if (!order) return <div>Order not found</div>

  return (
    <div>
      <h1>Order #{order.id}</h1>
      <p><strong>Status:</strong> {order.statusText}</p>
      <p><strong>Date:</strong> {new Date(order.createdOn).toLocaleDateString()}</p>
      <p><strong>Delivery Date:</strong> {new Date(order.deliveryDate).toLocaleDateString()}</p>
      {order.address && (
        <p><strong>Address:</strong> {order.address.addressLine1}, {order.address.city}, {order.address.state} {order.address.zipCode}</p>
      )}
      <h2>Items</h2>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '8px' }}>Book</th>
            <th style={{ textAlign: 'left', padding: '8px' }}>Price</th>
            <th style={{ textAlign: 'left', padding: '8px' }}>Qty</th>
          </tr>
        </thead>
        <tbody>
          {order.orderItems?.map(item => (
            <tr key={item.id}>
              <td style={{ padding: '8px' }}>{item.book.name}</td>
              <td style={{ padding: '8px' }}>${item.book.price.toFixed(2)}</td>
              <td style={{ padding: '8px' }}>{item.quantity}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <div style={{ marginTop: '15px' }}>
        <p>Subtotal: ${order.subTotal.toFixed(2)}</p>
        <p>Tax: ${order.tax.toFixed(2)}</p>
        <p><strong>Total: ${order.total.toFixed(2)}</strong></p>
      </div>
    </div>
  )
}
