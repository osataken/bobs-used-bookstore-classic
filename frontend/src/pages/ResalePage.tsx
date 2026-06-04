import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { getMyOffers } from '../services/api';
import { OfferStatus } from '../types';

const statusLabels: Record<number, string> = {
  [OfferStatus.PendingApproval]: 'Pending Approval',
  [OfferStatus.Approved]: 'Approved',
  [OfferStatus.Received]: 'Received',
  [OfferStatus.Paid]: 'Paid',
  [OfferStatus.Rejected]: 'Rejected',
};

export default function ResalePage() {
  const { data: offers } = useQuery({ queryKey: ['my-offers'], queryFn: getMyOffers });

  return (
    <div>
      <h1>My Resale Offers</h1>
      <Link to="/resale/new"><button style={{ marginBottom: '1rem' }}>Submit New Offer</button></Link>
      {!offers?.length ? (
        <p>No offers submitted yet.</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead><tr><th>Book</th><th>Author</th><th>Price</th><th>Status</th></tr></thead>
          <tbody>
            {offers.map(offer => (
              <tr key={offer.id} style={{ borderBottom: '1px solid #ddd' }}>
                <td style={{ padding: '0.5rem' }}>{offer.bookName}</td>
                <td>{offer.author}</td>
                <td>${offer.bookPrice.toFixed(2)}</td>
                <td>{statusLabels[offer.offerStatus]}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}
