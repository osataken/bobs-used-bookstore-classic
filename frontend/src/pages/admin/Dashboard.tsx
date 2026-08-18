import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { adminService } from '../../services';
import { useAuth } from '../../context/AuthContext';

export default function Dashboard() {
  const { user } = useAuth();
  const { data, isLoading } = useQuery({
    queryKey: ['admin-dashboard'],
    queryFn: () => adminService.getDashboard().then(r => r.data),
    enabled: !!user?.isAdmin,
  });

  if (!user?.isAdmin) return <p>Access denied.</p>;
  if (isLoading) return <p>Loading...</p>;

  return (
    <div>
      <h1 className="page-title">Admin Dashboard</h1>
      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
        <Link to="/admin/inventory"><button className="btn btn-primary">Inventory</button></Link>
        <Link to="/admin/orders"><button className="btn btn-primary">Orders</button></Link>
        <Link to="/admin/offers"><button className="btn btn-primary">Offers</button></Link>
        <Link to="/admin/reference-data"><button className="btn btn-primary">Reference Data</button></Link>
      </div>
      {data && (
        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '1rem' }}>
          <div className="card">
            <h3>Orders</h3>
            <p>Pending: {data.orderStatistics?.pending}</p>
            <p>Ordered: {data.orderStatistics?.ordered}</p>
            <p>Shipped: {data.orderStatistics?.shipped}</p>
            <p>Delivered: {data.orderStatistics?.delivered}</p>
          </div>
          <div className="card">
            <h3>Offers</h3>
            <p>Pending: {data.offerStatistics?.pendingApproval}</p>
            <p>Approved: {data.offerStatistics?.approved}</p>
            <p>Received: {data.offerStatistics?.received}</p>
            <p>Paid: {data.offerStatistics?.paid}</p>
          </div>
          <div className="card">
            <h3>Inventory</h3>
            <p>Total: {data.bookStatistics?.stockTotal}</p>
            <p>Low Stock: {data.bookStatistics?.lowStock}</p>
            <p>Out of Stock: {data.bookStatistics?.outOfStock}</p>
          </div>
        </div>
      )}
    </div>
  );
}
