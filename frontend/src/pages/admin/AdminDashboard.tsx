import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { getAdminDashboard } from '../../services/api';

export default function AdminDashboard() {
  const { data } = useQuery({ queryKey: ['admin-dashboard'], queryFn: getAdminDashboard });

  if (!data) return <div>Loading...</div>;

  return (
    <div>
      <h1>Admin Dashboard</h1>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '1rem' }}>
        <div style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
          <h2>Orders</h2>
          <p>Pending: {data.orders.pendingOrders}</p>
          <p>Past Due: {data.orders.pastDueOrders}</p>
          <p>This Month: {data.orders.ordersThisMonth}</p>
          <p>Total: {data.orders.ordersTotal}</p>
          <Link to="/admin/orders">Manage Orders →</Link>
        </div>
        <div style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
          <h2>Offers</h2>
          <p>Pending: {data.offers.pendingOffers}</p>
          <p>This Month: {data.offers.offersThisMonth}</p>
          <p>Total: {data.offers.offersTotal}</p>
          <Link to="/admin/offers">Manage Offers →</Link>
        </div>
        <div style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
          <h2>Inventory</h2>
          <p>Low Stock: {data.inventory.lowStock}</p>
          <p>Out of Stock: {data.inventory.outOfStock}</p>
          <p>Total: {data.inventory.stockTotal}</p>
          <Link to="/admin/inventory">Manage Inventory →</Link>
        </div>
      </div>
      <div style={{ marginTop: '1rem' }}>
        <Link to="/admin/reference-data">Manage Reference Data →</Link>
      </div>
    </div>
  );
}
