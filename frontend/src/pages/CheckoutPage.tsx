import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate, Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { getCart, getAddresses, checkout } from '../services/api';

export default function CheckoutPage() {
  const navigate = useNavigate();
  const { data: cart } = useQuery({ queryKey: ['cart'], queryFn: getCart });
  const { data: addresses } = useQuery({ queryKey: ['addresses'], queryFn: getAddresses });
  const [selectedAddressId, setSelectedAddressId] = useState<number>(0);

  const handleCheckout = async () => {
    if (!selectedAddressId) { toast.error('Please select an address'); return; }
    const result = await checkout(selectedAddressId);
    toast.success('Order placed successfully!');
    navigate(`/orders/${result.orderId}`);
  };

  const inStockItems = cart?.items?.filter(i => i.inStock) || [];

  return (
    <div>
      <h1>Checkout</h1>
      {!inStockItems.length ? (
        <p>No items in your cart are available for checkout.</p>
      ) : (
        <>
          <h2>Order Summary</h2>
          <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '1rem' }}>
            <thead><tr><th>Book</th><th>Price</th><th>Qty</th></tr></thead>
            <tbody>
              {inStockItems.map(item => (
                <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>{item.book.name}</td>
                  <td>${item.book.price.toFixed(2)}</td>
                  <td>{item.quantity}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <p><strong>Subtotal: ${cart?.subTotal.toFixed(2)}</strong></p>

          <h2>Delivery Address</h2>
          {!addresses?.length ? (
            <p>No addresses found. <Link to="/addresses/new?returnUrl=/checkout">Add an address</Link></p>
          ) : (
            <select value={selectedAddressId} onChange={e => setSelectedAddressId(Number(e.target.value))} style={{ padding: '0.5rem', width: '100%' }}>
              <option value={0}>Select an address...</option>
              {addresses.map(addr => (
                <option key={addr.id} value={addr.id}>
                  {addr.addressLine1}, {addr.city}, {addr.state} {addr.zipCode}
                </option>
              ))}
            </select>
          )}

          <button onClick={handleCheckout} style={{ marginTop: '1rem' }} disabled={!selectedAddressId}>
            Place Order
          </button>
        </>
      )}
    </div>
  );
}
