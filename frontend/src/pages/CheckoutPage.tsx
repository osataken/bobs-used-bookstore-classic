import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import api from '../services/api';
import { ShoppingCartItem, Address } from '../types';

function CheckoutPage() {
  const navigate = useNavigate();
  const [selectedAddressId, setSelectedAddressId] = useState<number>(0);

  const { data, isLoading } = useQuery({
    queryKey: ['checkout'],
    queryFn: async () => {
      const res = await api.get('/checkout');
      return res.data as {
        items: ShoppingCartItem[];
        subTotal: number;
        tax: number;
        total: number;
        addresses: Address[];
      };
    },
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedAddressId) {
      toast.error('Please select a delivery address');
      return;
    }
    try {
      const res = await api.post('/checkout', { selectedAddressId });
      toast.success('Order placed successfully!');
      navigate(`/orders/${res.data.orderId}`);
    } catch {
      toast.error('Failed to place order');
    }
  };

  if (isLoading) return <div>Loading...</div>;

  const items = data?.items || [];
  const addresses = data?.addresses || [];

  return (
    <div>
      <h1>Checkout</h1>
      {items.length === 0 ? (
        <p>Your cart is empty.</p>
      ) : (
        <form onSubmit={handleSubmit}>
          <h2>Order Summary</h2>
          <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '1rem' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Book</th>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Price</th>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Qty</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item) => (
                <tr key={item.id}>
                  <td style={{ padding: '0.5rem' }}>{item.book.name}</td>
                  <td style={{ padding: '0.5rem' }}>${item.book.price.toFixed(2)}</td>
                  <td style={{ padding: '0.5rem' }}>{item.quantity}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <p>Subtotal: ${data?.subTotal.toFixed(2)}</p>
          <p>Tax (10%): ${data?.tax.toFixed(2)}</p>
          <p><strong>Total: ${data?.total.toFixed(2)}</strong></p>

          <h2>Delivery Address</h2>
          {addresses.length === 0 ? (
            <p>No addresses found. Please add an address first.</p>
          ) : (
            <select value={selectedAddressId} onChange={(e) => setSelectedAddressId(Number(e.target.value))} required>
              <option value={0}>Select an address</option>
              {addresses.map((addr) => (
                <option key={addr.id} value={addr.id}>
                  {addr.addressLine1}, {addr.city}, {addr.state} {addr.zipCode}
                </option>
              ))}
            </select>
          )}

          <div style={{ marginTop: '1rem' }}>
            <button type="submit" disabled={addresses.length === 0}>Place Order</button>
          </div>
        </form>
      )}
    </div>
  );
}

export default CheckoutPage;
