import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation } from '@tanstack/react-query'
import { getCheckout, placeOrder } from '../services/api'
import toast from 'react-hot-toast'

function Checkout() {
  const navigate = useNavigate()
  const [selectedAddressId, setSelectedAddressId] = useState<number>(0)

  const { data, isLoading } = useQuery({
    queryKey: ['checkout'],
    queryFn: async () => {
      const res = await getCheckout()
      return res.data
    },
  })

  const placeOrderMutation = useMutation({
    mutationFn: () => placeOrder(selectedAddressId),
    onSuccess: (res) => {
      toast.success('Order placed successfully!')
      navigate(`/orders/${res.data.orderId}`)
    },
    onError: () => {
      toast.error('Failed to place order')
    },
  })

  if (isLoading) return <p>Loading...</p>
  if (!data) return <p>Error loading checkout</p>

  return (
    <div>
      <h1>Checkout</h1>
      <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '2rem' }}>
        <div>
          <h2>Order Items</h2>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left' }}>Book</th>
                <th style={{ textAlign: 'left' }}>Price</th>
              </tr>
            </thead>
            <tbody>
              {data.items?.map((item) => (
                <tr key={item.id}>
                  <td>{item.book?.name}</td>
                  <td>${item.book?.price.toFixed(2)}</td>
                </tr>
              ))}
            </tbody>
          </table>

          <h2>Delivery Address</h2>
          {data.addresses.length === 0 ? (
            <p>No addresses found. Please add an address first.</p>
          ) : (
            <select value={selectedAddressId} onChange={e => setSelectedAddressId(Number(e.target.value))} style={{ width: '100%', padding: '0.5rem' }}>
              <option value={0}>Select an address</option>
              {data.addresses.map(addr => (
                <option key={addr.id} value={addr.id}>
                  {addr.addressLine1}, {addr.city}, {addr.state} {addr.zipCode}
                </option>
              ))}
            </select>
          )}
        </div>
        <div style={{ background: '#f9f9f9', padding: '1rem', borderRadius: '8px' }}>
          <h2>Order Summary</h2>
          <p>Subtotal: ${data.subTotal.toFixed(2)}</p>
          <p>Tax (10%): ${data.tax.toFixed(2)}</p>
          <p style={{ fontWeight: 'bold' }}>Total: ${data.total.toFixed(2)}</p>
          <button
            onClick={() => placeOrderMutation.mutate()}
            disabled={selectedAddressId === 0 || placeOrderMutation.isPending}
            style={{ width: '100%', padding: '0.75rem', marginTop: '1rem' }}
          >
            Place Order
          </button>
        </div>
      </div>
    </div>
  )
}

export default Checkout
