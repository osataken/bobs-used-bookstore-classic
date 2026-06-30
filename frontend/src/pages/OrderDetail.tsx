import { useParams } from 'react-router-dom'
import { useOrderDetails } from '../hooks/useOrders'

function OrderDetail() {
  const { id } = useParams<{ id: string }>()
  const { data, isLoading } = useOrderDetails(Number(id))

  if (isLoading) return <p>Loading...</p>
  if (!data) return <p>Order not found</p>

  const { order, subTotal, tax, total } = data

  return (
    <div>
      <h1>Order #{order.id}</h1>
      <p><strong>Status:</strong> {order.orderStatus}</p>
      <p><strong>Delivery Date:</strong> {order.deliveryDate}</p>
      {order.address && (
        <p><strong>Address:</strong> {order.address.addressLine1}, {order.address.city}, {order.address.state} {order.address.zipCode}</p>
      )}
      <h2>Items</h2>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left' }}>Book</th>
            <th style={{ textAlign: 'left' }}>Quantity</th>
            <th style={{ textAlign: 'left' }}>Price</th>
          </tr>
        </thead>
        <tbody>
          {order.orderItems?.map(item => (
            <tr key={item.id}>
              <td>{item.book?.name}</td>
              <td>{item.quantity}</td>
              <td>${item.book?.price.toFixed(2)}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <div style={{ marginTop: '1rem' }}>
        <p>Subtotal: ${subTotal.toFixed(2)}</p>
        <p>Tax: ${tax.toFixed(2)}</p>
        <p><strong>Total: ${total.toFixed(2)}</strong></p>
      </div>
    </div>
  )
}

export default OrderDetail
