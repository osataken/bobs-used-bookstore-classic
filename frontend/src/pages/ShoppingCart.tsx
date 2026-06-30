import { useCart, useDeleteCartItem } from '../hooks/useCart'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

function ShoppingCart() {
  const { data, isLoading } = useCart()
  const deleteItem = useDeleteCartItem()
  const { user } = useAuth()

  if (isLoading) return <p>Loading...</p>

  const items = data?.items || []

  return (
    <div>
      <h1>Shopping Cart</h1>
      {items.length === 0 ? (
        <p>Your cart is empty. <Link to="/search">Browse books</Link></p>
      ) : (
        <>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Book</th>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Price</th>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Quantity</th>
                <th style={{ padding: '0.5rem' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item) => (
                <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>{item.book?.name}</td>
                  <td style={{ padding: '0.5rem' }}>${item.book?.price.toFixed(2)}</td>
                  <td style={{ padding: '0.5rem' }}>{item.quantity}</td>
                  <td style={{ padding: '0.5rem' }}>
                    <button onClick={() => deleteItem.mutate(item.id)}>Remove</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <div style={{ marginTop: '1rem', textAlign: 'right' }}>
            <p><strong>Subtotal: ${data?.subTotal.toFixed(2)}</strong></p>
            {user ? (
              <Link to="/checkout"><button>Proceed to Checkout</button></Link>
            ) : (
              <p>Please log in to checkout</p>
            )}
          </div>
        </>
      )}
    </div>
  )
}

export default ShoppingCart
