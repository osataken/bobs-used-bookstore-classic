import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import api from '../../services/api';

function DashboardPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['adminDashboard'],
    queryFn: async () => {
      const res = await api.get('/admin/dashboard');
      return res.data;
    },
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Admin Dashboard</h1>
      <nav style={{ marginBottom: '1rem', display: 'flex', gap: '1rem' }}>
        <Link to="/admin/inventory">Inventory</Link>
        <Link to="/admin/orders">Orders</Link>
        <Link to="/admin/offers">Offers</Link>
        <Link to="/admin/reference-data">Reference Data</Link>
      </nav>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '1rem' }}>
        <div style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
          <h2>Orders</h2>
          <p>Pending: {data?.orders.pendingOrders}</p>
          <p>Past Due: {data?.orders.pastDueOrders}</p>
          <p>This Month: {data?.orders.ordersThisMonth}</p>
          <p>Total: {data?.orders.ordersTotal}</p>
        </div>
        <div style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
          <h2>Offers</h2>
          <p>Pending: {data?.offers.pendingOffers}</p>
          <p>This Month: {data?.offers.offersThisMonth}</p>
          <p>Total: {data?.offers.offersTotal}</p>
        </div>
        <div style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
          <h2>Inventory</h2>
          <p>Low Stock: {data?.inventory.lowStock}</p>
          <p>Out of Stock: {data?.inventory.outOfStock}</p>
          <p>Total: {data?.inventory.stockTotal}</p>
        </div>
      </div>
    </div>
  );
}

export default DashboardPage;
