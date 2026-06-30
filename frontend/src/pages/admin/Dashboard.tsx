import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getAdminDashboard } from '../../services/api'

function Dashboard() {
  const { data, isLoading } = useQuery({
    queryKey: ['adminDashboard'],
    queryFn: async () => { const res = await getAdminDashboard(); return res.data; },
  })

  if (isLoading) return <p>Loading...</p>
  if (!data) return <p>Error loading dashboard</p>

  return (
    <div>
      <h1>Admin Dashboard</h1>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))', gap: '1rem' }}>
        <div style={{ padding: '1rem', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h3>Orders</h3>
          <p>Total: {data.totalOrders}</p>
          <p>Pending: {data.pendingOrders}</p>
          <Link to="/admin/orders">Manage Orders</Link>
        </div>
        <div style={{ padding: '1rem', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h3>Offers</h3>
          <p>Pending: {data.pendingOffers}</p>
          <Link to="/admin/offers">Manage Offers</Link>
        </div>
        <div style={{ padding: '1rem', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h3>Inventory</h3>
          <p>Total Books: {data.totalBooks}</p>
          <p>Low Stock: {data.lowStockBooks}</p>
          <Link to="/admin/inventory">Manage Inventory</Link>
        </div>
        <div style={{ padding: '1rem', border: '1px solid #ddd', borderRadius: '8px' }}>
          <h3>Reference Data</h3>
          <Link to="/admin/reference-data">Manage Reference Data</Link>
        </div>
      </div>
    </div>
  )
}

export default Dashboard
