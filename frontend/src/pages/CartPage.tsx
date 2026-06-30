import { useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getCart, deleteCartItem } from '../services/api'
import toast from 'react-hot-toast'

export default function CartPage() {
  const queryClient = useQueryClient()
  const { data, isLoading } = useQuery({
    queryKey: ['cart'],
    queryFn: () => getCart().then(r => r.data),
  })

  const handleDelete = async (itemId: number) => {
    try {
      await deleteCartItem(itemId)
      queryClient.invalidateQueries({ queryKey: ['cart'] })
      toast.success('Item removed')
    } catch {
      toast.error('Failed to remove item')
    }
  }

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      <h1>Shopping Cart</h1>
      {(!data?.items || data.items.length === 0) ? (
        <p>Your cart is empty. <Link to="/search">Browse books</Link></p>
      ) : (
        <>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Book</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Price</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Quantity</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {data.items.map(item => (
                <tr key={item.id}>
                  <td style={{ padding: '8px' }}><Link to={`/search/${item.bookId}`}>{item.book.name}</Link></td>
                  <td style={{ padding: '8px' }}>${item.book.price.toFixed(2)}</td>
                  <td style={{ padding: '8px' }}>{item.quantity}</td>
                  <td style={{ padding: '8px' }}><button onClick={() => handleDelete(item.id)}>Remove</button></td>
                </tr>
              ))}
            </tbody>
          </table>
          <div style={{ marginTop: '20px', textAlign: 'right' }}>
            <p><strong>Subtotal: ${data.subTotal.toFixed(2)}</strong></p>
            <Link to="/checkout"><button>Proceed to Checkout</button></Link>
          </div>
        </>
      )}
    </div>
  )
}
