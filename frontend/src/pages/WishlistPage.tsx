import { useQuery, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import api from '../services/api';
import { ShoppingCartItem } from '../types';

function WishlistPage() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['wishlist'],
    queryFn: async () => {
      const res = await api.get('/wishlist');
      return res.data as { items: ShoppingCartItem[] };
    },
  });

  const moveToCart = async (itemId: number) => {
    try {
      await api.post('/wishlist/move', { shoppingCartItemId: itemId });
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      toast.success('Moved to cart');
    } catch {
      toast.error('Failed to move item');
    }
  };

  const moveAllToCart = async () => {
    try {
      await api.post('/wishlist/moveAll');
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      toast.success('All items moved to cart');
    } catch {
      toast.error('Failed to move items');
    }
  };

  const removeItem = async (itemId: number) => {
    try {
      await api.post('/wishlist/delete', { shoppingCartItemId: itemId });
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      toast.success('Item removed');
    } catch {
      toast.error('Failed to remove item');
    }
  };

  if (isLoading) return <div>Loading...</div>;

  const items = data?.items || [];

  return (
    <div>
      <h1>Wishlist</h1>
      {items.length === 0 ? (
        <p>Your wishlist is empty.</p>
      ) : (
        <>
          <button onClick={moveAllToCart} style={{ marginBottom: '1rem' }}>Move All to Cart</button>
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
                  <td style={{ padding: '0.5rem' }}>{item.book.name}</td>
                  <td style={{ padding: '0.5rem' }}>${item.book.price.toFixed(2)}</td>
                  <td style={{ padding: '0.5rem', display: 'flex', gap: '0.5rem' }}>
                    <button onClick={() => moveToCart(item.id)}>Move to Cart</button>
                    <button onClick={() => removeItem(item.id)}>Remove</button>
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

export default WishlistPage;
