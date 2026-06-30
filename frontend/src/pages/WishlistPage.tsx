import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getWishlist, moveToCart, moveAllToCart, deleteWishlistItem } from '../services/api'
import toast from 'react-hot-toast'

export default function WishlistPage() {
  const queryClient = useQueryClient()
  const { data, isLoading } = useQuery({
    queryKey: ['wishlist'],
    queryFn: () => getWishlist().then(r => r.data),
  })

  const handleMove = async (itemId: number) => {
    try {
      await moveToCart(itemId)
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success('Moved to cart')
    } catch {
      toast.error('Failed to move item')
    }
  }

  const handleMoveAll = async () => {
    try {
      await moveAllToCart()
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success('All items moved to cart')
    } catch {
      toast.error('Failed to move items')
    }
  }

  const handleDelete = async (itemId: number) => {
    try {
      await deleteWishlistItem(itemId)
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success('Item removed')
    } catch {
      toast.error('Failed to remove item')
    }
  }

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      <h1>Wishlist</h1>
      {(!data?.items || data.items.length === 0) ? (
        <p>Your wishlist is empty.</p>
      ) : (
        <>
          <button onClick={handleMoveAll} style={{ marginBottom: '15px' }}>Move All to Cart</button>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Book</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Price</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {data.items.map(item => (
                <tr key={item.id}>
                  <td style={{ padding: '8px' }}>{item.book.name}</td>
                  <td style={{ padding: '8px' }}>${item.book.price.toFixed(2)}</td>
                  <td style={{ padding: '8px' }}>
                    <button onClick={() => handleMove(item.id)}>Move to Cart</button>{' '}
                    <button onClick={() => handleDelete(item.id)}>Remove</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </div>
  )
}
