import { useState } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import api from '../../services/api';
import { Order, PaginatedResponse, OrderStatusLabels } from '../../types';
import Pagination from '../../components/Pagination';

function AdminOrdersPage() {
  const queryClient = useQueryClient();
  const [pageIndex, setPageIndex] = useState(1);
  const [statusFilter, setStatusFilter] = useState('');

  const { data, isLoading } = useQuery({
    queryKey: ['adminOrders', pageIndex, statusFilter],
    queryFn: async () => {
      const params: Record<string, string | number> = { pageIndex, pageSize: 10 };
      if (statusFilter) params.status = statusFilter;
      const res = await api.get('/admin/orders', { params });
      return res.data as PaginatedResponse<Order>;
    },
  });

  const updateStatus = async (orderId: number, orderStatus: number) => {
    try {
      await api.post('/admin/orders/status', { orderId, orderStatus });
      toast.success('Order status updated');
      queryClient.invalidateQueries({ queryKey: ['adminOrders'] });
    } catch {
      toast.error('Failed to update status');
    }
  };

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Order Management</h1>
      <select value={statusFilter} onChange={(e) => { setStatusFilter(e.target.value); setPageIndex(1); }}>
        <option value="">All Statuses</option>
        {Object.entries(OrderStatusLabels).map(([val, label]) => (
          <option key={val} value={val}>{label}</option>
        ))}
      </select>

      <table style={{ width: '100%', borderCollapse: 'collapse', marginTop: '1rem' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>ID</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Status</th>
            <th style={{ textAlign: 'left', padding: '0.5rem' }}>Delivery</th>
            <th style={{ padding: '0.5rem' }}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {data?.items.map((order) => (
            <tr key={order.id} style={{ borderBottom: '1px solid #ddd' }}>
              <td style={{ padding: '0.5rem' }}>#{order.id}</td>
              <td style={{ padding: '0.5rem' }}>{OrderStatusLabels[order.orderStatus]}</td>
              <td style={{ padding: '0.5rem' }}>{new Date(order.deliveryDate).toLocaleDateString()}</td>
              <td style={{ padding: '0.5rem' }}>
                <select onChange={(e) => updateStatus(order.id, Number(e.target.value))} defaultValue="">
                  <option value="" disabled>Change Status</option>
                  {Object.entries(OrderStatusLabels).map(([val, label]) => (
                    <option key={val} value={val}>{label}</option>
                  ))}
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
    </div>
  );
}

export default AdminOrdersPage;
