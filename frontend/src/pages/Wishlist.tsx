import { useQuery, useQueryClient } from '@tanstack/react-query';
import { wishlistService } from '../services';
import toast from 'react-hot-toast';

export default function Wishlist() {
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ['wishlist'],
    queryFn: () => wishlistService.getWishlist().then(r => r.data),
  });

  const moveToCart = async (itemId: number) => {
    try {
      await wishlistService.moveToCart(itemId);
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      toast.success('Moved to cart');
    } catch { toast.error('Failed'); }
  };

  const moveAll = async () => {
    try {
      await wishlistService.moveAllToCart();
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      toast.success('All items moved to cart');
    } catch { toast.error('Failed'); }
  };

  const remove = async (itemId: number) => {
    try {
      await wishlistService.deleteItem(itemId);
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      toast.success('Removed');
    } catch { toast.error('Failed'); }
  };

  if (isLoading) return <p>Loading...</p>;

  return (
    <div>
      <h1 className="page-title">Wishlist</h1>
      {!data?.items?.length ? <p>Your wishlist is empty.</p> : (
        <>
          <button className="btn btn-primary" onClick={moveAll} style={{ marginBottom: '1rem' }}>Move All to Cart</button>
          <table>
            <thead><tr><th>Book</th><th>Price</th><th>Actions</th></tr></thead>
            <tbody>
              {data.items.map(item => (
                <tr key={item.id}>
                  <td>{item.book?.name}</td>
                  <td>${item.book?.price.toFixed(2)}</td>
                  <td>
                    <button className="btn btn-primary" onClick={() => moveToCart(item.id)}>Move to Cart</button>
                    <button className="btn btn-danger" onClick={() => remove(item.id)} style={{ marginLeft: 4 }}>Remove</button>
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
