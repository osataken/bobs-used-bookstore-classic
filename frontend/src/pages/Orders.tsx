import { Link } from 'react-router-dom'
import { useOrders, useCancelOrder } from '../hooks/useOrders'
import { OrderStatus } from '../types'

function Orders() {
  const { data: orders, isLoading } = useOrders()
  const cancelOrder = useCancelOrder()

  if (isLoading) return <p>Loading...</p>

  return (
    <div>
      <h1>My Orders</h1>
      {(!orders || orders.length === 0) ? (
        <p>No orders found.</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Order #</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Date</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Status</th>
              <th style={{ textAlign: 'left', padding: '0.5rem' }}>Total</th>
              <th style={{ padding: '0.5rem' }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {orders.map(order => (
              <tr key={order.id} style={{ borderBottom: '1px solid #ddd' }}>
                <td style={{ padding: '0.5rem' }}><Link to={`/orders/${order.id}`}>#{order.id}</Link></td>
                <td style={{ padding: '0.5rem' }}>{order.createdOn}</td>
                <td style={{ padding: '0.5rem' }}>{order.statusText}</td>
                <td style={{ padding: '0.5rem' }}>${order.total?.toFixed(2)}</td>
                <td style={{ padding: '0.5rem' }}>
                  {order.orderStatus === OrderStatus.Pending && (
                    <button onClick={() => cancelOrder.mutate(order.id)}>Cancel</button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

export default Orders
