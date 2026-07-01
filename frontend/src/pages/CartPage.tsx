import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import api from '../services/api';
import { ShoppingCartItem } from '../types';

function CartPage() {
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['cart'],
    queryFn: async () => {
      const res = await api.get('/cart');
      return res.data as { items: ShoppingCartItem[]; subTotal: number };
    },
  });

  const removeItem = async (itemId: number) => {
    try {
      await api.post('/cart/delete', { shoppingCartItemId: itemId });
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      toast.success('Item removed');
    } catch {
      toast.error('Failed to remove item');
    }
  };

  if (isLoading) return <div>Loading...</div>;

  const items = data?.items || [];

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
                  <td style={{ padding: '0.5rem' }}>{item.book.name}</td>
                  <td style={{ padding: '0.5rem' }}>${item.book.price.toFixed(2)}</td>
                  <td style={{ padding: '0.5rem' }}>{item.quantity}</td>
                  <td style={{ padding: '0.5rem' }}>
                    <button onClick={() => removeItem(item.id)}>Remove</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <p style={{ fontWeight: 'bold', marginTop: '1rem' }}>Subtotal: ${data?.subTotal.toFixed(2)}</p>
          <Link to="/checkout"><button>Proceed to Checkout</button></Link>
        </>
      )}
    </div>
  );
}

export default CartPage;
