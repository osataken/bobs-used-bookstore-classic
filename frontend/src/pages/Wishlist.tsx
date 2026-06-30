import { useWishlist, useMoveToCart, useMoveAllToCart, useDeleteWishlistItem } from '../hooks/useCart'

function Wishlist() {
  const { data, isLoading } = useWishlist()
  const moveToCart = useMoveToCart()
  const moveAll = useMoveAllToCart()
  const deleteItem = useDeleteWishlistItem()

  if (isLoading) return <p>Loading...</p>

  const items = data?.items || []

  return (
    <div>
      <h1>Wishlist</h1>
      {items.length === 0 ? (
        <p>Your wishlist is empty.</p>
      ) : (
        <>
          <button onClick={() => moveAll.mutate()} style={{ marginBottom: '1rem' }}>Move All to Cart</button>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Book</th>
                <th style={{ textAlign: 'left', padding: '0.5rem' }}>Price</th>
                <th style={{ padding: '0.5rem' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {items.map((item) => (
                <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>{item.book?.name}</td>
                  <td style={{ padding: '0.5rem' }}>${item.book?.price.toFixed(2)}</td>
                  <td style={{ padding: '0.5rem', display: 'flex', gap: '0.5rem' }}>
                    <button onClick={() => moveToCart.mutate(item.id)}>Move to Cart</button>
                    <button onClick={() => deleteItem.mutate(item.id)}>Remove</button>
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

export default Wishlist
