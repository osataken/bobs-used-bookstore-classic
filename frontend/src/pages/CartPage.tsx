import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { getCart, deleteCartItem } from '../services/api';
import { useAuth } from '../context/AuthContext';

export default function CartPage() {
  const queryClient = useQueryClient();
  const { auth } = useAuth();
  const { data } = useQuery({ queryKey: ['cart'], queryFn: getCart });

  const handleDelete = async (itemId: number) => {
    await deleteCartItem(itemId);
    toast.success('Item removed from shopping cart');
    queryClient.invalidateQueries({ queryKey: ['cart'] });
  };

  return (
    <div>
      <h1>Shopping Cart</h1>
      {!data?.items?.length ? (
        <p>Your cart is empty. <Link to="/search">Browse books</Link></p>
      ) : (
        <>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr><th>Book</th><th>Price</th><th>Quantity</th><th>Status</th><th>Action</th></tr>
            </thead>
            <tbody>
              {data.items.map(item => (
                <tr key={item.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>{item.book.name}</td>
                  <td>${item.book.price.toFixed(2)}</td>
                  <td>{item.quantity}</td>
                  <td>{item.inStock ? 'In Stock' : <span style={{ color: 'red' }}>Out of Stock</span>}</td>
                  <td><button onClick={() => handleDelete(item.id)}>Remove</button></td>
                </tr>
              ))}
            </tbody>
          </table>
          <p style={{ marginTop: '1rem' }}><strong>Subtotal: ${data.subTotal.toFixed(2)}</strong></p>
          {auth.authenticated && (
            <Link to="/checkout"><button style={{ marginTop: '1rem' }}>Proceed to Checkout</button></Link>
          )}
        </>
      )}
    </div>
  );
}
