import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getAdminOrders, updateOrderStatus } from '../../services/api'
import Pagination from '../../components/Pagination'
import toast from 'react-hot-toast'

export default function AdminOrdersPage() {
  const queryClient = useQueryClient()
  const [pageIndex, setPageIndex] = useState(1)

  const { data, isLoading } = useQuery({
    queryKey: ['adminOrders', pageIndex],
    queryFn: () => getAdminOrders({ pageIndex, pageSize: 10 }).then(r => r.data),
  })

  const handleStatusUpdate = async (id: number, status: number) => {
    try {
      await updateOrderStatus(id, status)
      queryClient.invalidateQueries({ queryKey: ['adminOrders'] })
      toast.success('Status updated')
    } catch {
      toast.error('Failed to update status')
    }
  }

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      <h1>Order Management</h1>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>ID</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Customer</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Status</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Total</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {data?.items?.map(order => (
            <tr key={order.id}>
              <td style={{ padding: '8px' }}>#{order.id}</td>
              <td style={{ padding: '8px' }}>{order.customerName}</td>
              <td style={{ padding: '8px' }}>{order.statusText}</td>
              <td style={{ padding: '8px' }}>${order.total.toFixed(2)}</td>
              <td style={{ padding: '8px' }}>
                <select value={order.orderStatus} onChange={e => handleStatusUpdate(order.id, Number(e.target.value))}>
                  <option value={0}>Pending</option>
                  <option value={1}>Ordered</option>
                  <option value={2}>Shipped</option>
                  <option value={3}>Delivered</option>
                  <option value={4}>Cancelled</option>
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
    </div>
  )
}
