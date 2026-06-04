import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { getWishlist, moveToCart, moveAllToCart, deleteWishlistItem } from '../services/api';

export default function WishlistPage() {
  const queryClient = useQueryClient();
  const { data } = useQuery({ queryKey: ['wishlist'], queryFn: getWishlist });

  const handleMove = async (itemId: number) => {
    await moveToCart(itemId);
    toast.success('Item moved to shopping cart');
    queryClient.invalidateQueries({ queryKey: ['wishlist'] });
    queryClient.invalidateQueries({ queryKey: ['cart'] });
  };

  const handleMoveAll = async () => {
    await moveAllToCart();
    toast.success('All items moved to shopping cart');
    queryClient.invalidateQueries({ queryKey: ['wishlist'] });
    queryClient.invalidateQueries({ queryKey: ['cart'] });
  };

  const handleDelete = async (itemId: number) => {
    await deleteWishlistItem(itemId);
    toast.success('Item removed from wishlist');
    queryClient.invalidateQueries({ queryKey: ['wishlist'] });
  };

  return (
    <div>
      <h1>Wishlist</h1>
      {!data?.items?.length ? (
        <p>Your wishlist is empty. <Link to="/search">Browse books</Link></p>
      ) : (
        <>
          <button onClick={handleMoveAll} style={{ marginBottom: '1rem' }}>Move All to Cart</button>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead><tr><th>Book</th><th>Price</th><th>Actions</th></tr></thead>
            <tbody>
              {data.items.map(item => (
                <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>{item.book.name}</td>
                  <td>${item.book.price.toFixed(2)}</td>
                  <td>
                    <button onClick={() => handleMove(item.id)} style={{ marginRight: '0.5rem' }}>Move to Cart</button>
                    <button onClick={() => handleDelete(item.id)}>Remove</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </>
      )}
    </div>
  );
}
