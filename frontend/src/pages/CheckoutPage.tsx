import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate, Link } from 'react-router-dom'
import { getCheckout, submitCheckout } from '../services/api'
import toast from 'react-hot-toast'

export default function CheckoutPage() {
  const navigate = useNavigate()
  const [selectedAddressId, setSelectedAddressId] = useState<number | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const { data, isLoading } = useQuery({
    queryKey: ['checkout'],
    queryFn: () => getCheckout().then(r => r.data),
  })

  const handleSubmit = async () => {
    if (!selectedAddressId) {
      toast.error('Please select a delivery address')
      return
    }
    setSubmitting(true)
    try {
      const response = await submitCheckout(selectedAddressId)
      toast.success('Order placed successfully!')
      navigate(`/orders/${response.data.orderId}`)
    } catch {
      toast.error('Failed to place order')
    } finally {
      setSubmitting(false)
    }
  }

  if (isLoading) return <div>Loading...</div>

  const cart = data?.cart
  const addresses = data?.addresses

  if (!cart?.items || cart.items.length === 0) {
    return <div><h1>Checkout</h1><p>Your cart is empty.</p></div>
  }

  const tax = cart.subTotal * 0.1
  const total = cart.subTotal + tax

  return (
    <div>
      <h1>Checkout</h1>
      <h2>Order Summary</h2>
      <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '20px' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '8px' }}>Book</th>
            <th style={{ textAlign: 'left', padding: '8px' }}>Price</th>
            <th style={{ textAlign: 'left', padding: '8px' }}>Qty</th>
          </tr>
        </thead>
        <tbody>
          {cart.items.map(item => (
            <tr key={item.id}>
              <td style={{ padding: '8px' }}>{item.book.name}</td>
              <td style={{ padding: '8px' }}>${item.book.price.toFixed(2)}</td>
              <td style={{ padding: '8px' }}>{item.quantity}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <p>Subtotal: ${cart.subTotal.toFixed(2)}</p>
      <p>Tax (10%): ${tax.toFixed(2)}</p>
      <p><strong>Total: ${total.toFixed(2)}</strong></p>

      <h2>Delivery Address</h2>
      {(!addresses || addresses.length === 0) ? (
        <p>No addresses found. <Link to="/addresses">Add an address</Link></p>
      ) : (
        <div>
          {addresses.map(addr => (
            <label key={addr.id} style={{ display: 'block', padding: '10px', border: '1px solid #ddd', marginBottom: '5px', cursor: 'pointer' }}>
              <input type="radio" name="address" checked={selectedAddressId === addr.id} onChange={() => setSelectedAddressId(addr.id)} />
              {' '}{addr.addressLine1}, {addr.city}, {addr.state} {addr.zipCode}, {addr.country}
            </label>
          ))}
        </div>
      )}

      <button onClick={handleSubmit} disabled={submitting || !selectedAddressId} style={{ marginTop: '20px', padding: '10px 20px' }}>
        {submitting ? 'Placing Order...' : 'Place Order'}
      </button>
    </div>
  )
}
