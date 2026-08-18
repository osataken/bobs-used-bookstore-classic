import { useQuery, useQueryClient } from '@tanstack/react-query';
import { cartService } from '../services';
import toast from 'react-hot-toast';

export default function Cart() {
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ['cart'],
    queryFn: () => cartService.getCart().then(r => r.data),
  });

  const removeItem = async (itemId: number) => {
    try {
      await cartService.deleteItem(itemId);
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      toast.success('Item removed');
    } catch { toast.error('Failed to remove item'); }
  };

  if (isLoading) return <p>Loading...</p>;

  return (
    <div>
      <h1 className="page-title">Shopping Cart</h1>
      {!data?.items?.length ? <p>Your cart is empty.</p> : (
        <>
          <table>
            <thead><tr><th>Book</th><th>Price</th><th>Quantity</th><th>Actions</th></tr></thead>
            <tbody>
              {data.items.map(item => (
                <tr key={item.id}>
                  <td>{item.book?.name}</td>
                  <td>${item.book?.price.toFixed(2)}</td>
                  <td>{item.quantity}</td>
                  <td><button className="btn btn-danger" onClick={() => removeItem(item.id)}>Remove</button></td>
                </tr>
              ))}
            </tbody>
          </table>
          <p style={{ marginTop: '1rem' }}><strong>Subtotal: ${data.subTotal.toFixed(2)}</strong></p>
          <a href="/checkout"><button className="btn btn-primary" style={{ marginTop: '1rem' }}>Proceed to Checkout</button></a>
        </>
      )}
    </div>
  );
}
