import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { checkoutService } from '../services';
import { Address } from '../types';
import { useAuth } from '../context/AuthContext';
import toast from 'react-hot-toast';

export default function Checkout() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [selectedAddress, setSelectedAddress] = useState<number>(0);
  const [checkoutData, setCheckoutData] = useState<{ subTotal: number; tax: number; total: number } | null>(null);

  useEffect(() => {
    if (!user) return;
    checkoutService.getCheckout().then(r => {
      setAddresses(r.data.addresses || []);
      setCheckoutData({ subTotal: r.data.subTotal, tax: r.data.tax, total: r.data.total });
      if (r.data.addresses?.length > 0) setSelectedAddress(r.data.addresses[0].id);
    });
  }, [user]);

  const submit = async () => {
    if (!selectedAddress) { toast.error('Please select an address'); return; }
    try {
      const res = await checkoutService.submitOrder(selectedAddress);
      toast.success('Order placed!');
      navigate(`/orders/${res.data.orderId}`);
    } catch (e: unknown) {
      const err = e as { response?: { data?: { error?: string } } };
      toast.error(err.response?.data?.error || 'Failed to place order');
    }
  };

  if (!user) return <p>Please log in to checkout.</p>;

  return (
    <div>
      <h1 className="page-title">Checkout</h1>
      {checkoutData && (
        <div className="card">
          <p>Subtotal: ${checkoutData.subTotal.toFixed(2)}</p>
          <p>Tax (10%): ${checkoutData.tax.toFixed(2)}</p>
          <p><strong>Total: ${checkoutData.total.toFixed(2)}</strong></p>
        </div>
      )}
      <h2>Select Delivery Address</h2>
      {addresses.length === 0 ? <p>No addresses. <a href="/addresses">Add one</a></p> : (
        <select value={selectedAddress} onChange={e => setSelectedAddress(Number(e.target.value))}>
          {addresses.map(a => <option key={a.id} value={a.id}>{a.addressLine1}, {a.city}, {a.state} {a.zipCode}</option>)}
        </select>
      )}
      <button className="btn btn-primary" onClick={submit} style={{ marginTop: '1rem' }}>Place Order</button>
    </div>
  );
}
