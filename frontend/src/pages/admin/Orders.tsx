import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getAdminOrders, updateOrderStatus } from '../../services/api'
import toast from 'react-hot-toast'

function AdminOrders() {
  const queryClient = useQueryClient()
  const [pageIndex, setPageIndex] = useState(1)

  const { data } = useQuery({
    queryKey: ['adminOrders', pageIndex],
    queryFn: async () => { const res = await getAdminOrders({ pageIndex, pageSize: 10 }); return res.data; },
  })

  const updateStatus = useMutation({
    mutationFn: ({ orderId, status }: { orderId: number; status: number }) => updateOrderStatus(orderId, status),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminOrders'] }); toast.success('Status updated'); },
  })

  const statusLabels = ['Pending', 'Ordered', 'Shipped', 'Delivered', 'Cancelled']

  return (
    <div>
      <h1>Order Management</h1>
      {data && (
        <>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead><tr><th style={{ textAlign: 'left' }}>Order #</th><th style={{ textAlign: 'left' }}>Customer</th><th style={{ textAlign: 'left' }}>Status</th><th>Update Status</th></tr></thead>
            <tbody>
              {data.items.map(order => (
                <tr key={order.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>#{order.id}</td>
                  <td style={{ padding: '0.5rem' }}>{order.customer?.firstName} {order.customer?.lastName}</td>
                  <td style={{ padding: '0.5rem' }}>{statusLabels[order.orderStatus]}</td>
                  <td style={{ padding: '0.5rem' }}>
                    <select value={order.orderStatus} onChange={e => updateStatus.mutate({ orderId: order.id, status: Number(e.target.value) })}>
                      {statusLabels.map((label, i) => <option key={i} value={i}>{label}</option>)}
                    </select>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {data.totalPages > 1 && (
            <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
              <button onClick={() => setPageIndex(p => Math.max(1, p - 1))} disabled={pageIndex === 1}>Previous</button>
              <span>Page {pageIndex} of {data.totalPages}</span>
              <button onClick={() => setPageIndex(p => p + 1)} disabled={pageIndex === data.totalPages}>Next</button>
            </div>
          )}
        </>
      )}
    </div>
  )
}

export default AdminOrders
