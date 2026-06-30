import { useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getOrders, cancelOrder } from '../services/api'
import toast from 'react-hot-toast'

export default function OrdersPage() {
  const queryClient = useQueryClient()
  const { data, isLoading } = useQuery({
    queryKey: ['orders'],
    queryFn: () => getOrders().then(r => r.data),
  })

  const handleCancel = async (id: number) => {
    if (!confirm('Are you sure you want to cancel this order?')) return
    try {
      await cancelOrder(id)
      queryClient.invalidateQueries({ queryKey: ['orders'] })
      toast.success('Order cancelled')
    } catch {
      toast.error('Failed to cancel order')
    }
  }

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      <h1>My Orders</h1>
      {(!data?.orders || data.orders.length === 0) ? (
        <p>No orders yet.</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Order #</th>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Date</th>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Status</th>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Total</th>
              <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {data.orders.map(order => (
              <tr key={order.id}>
                <td style={{ padding: '8px' }}><Link to={`/orders/${order.id}`}>#{order.id}</Link></td>
                <td style={{ padding: '8px' }}>{new Date(order.createdOn).toLocaleDateString()}</td>
                <td style={{ padding: '8px' }}>{order.statusText}</td>
                <td style={{ padding: '8px' }}>${order.total.toFixed(2)}</td>
                <td style={{ padding: '8px' }}>
                  {order.orderStatus === 0 && <button onClick={() => handleCancel(order.id)}>Cancel</button>}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
