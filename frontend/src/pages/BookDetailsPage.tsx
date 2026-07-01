import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import api from '../services/api';
import { Book } from '../types';

function BookDetailsPage() {
  const { id } = useParams<{ id: string }>();

  const { data: book, isLoading } = useQuery({
    queryKey: ['book', id],
    queryFn: async () => {
      const res = await api.get(`/search/${id}`);
      return res.data as Book;
    },
  });

  const addToCart = async () => {
    try {
      await api.post('/cart/add', { bookId: Number(id), quantity: 1 });
      toast.success('Added to cart');
    } catch {
      toast.error('Failed to add to cart');
    }
  };

  const addToWishlist = async () => {
    try {
      await api.post('/wishlist/add', { bookId: Number(id) });
      toast.success('Added to wishlist');
    } catch {
      toast.error('Failed to add to wishlist');
    }
  };

  if (isLoading) return <div>Loading...</div>;
  if (!book) return <div>Book not found</div>;

  return (
    <div>
      <h1>{book.name}</h1>
      <p><strong>Author:</strong> {book.author}</p>
      <p><strong>Price:</strong> ${book.price.toFixed(2)}</p>
      <p><strong>ISBN:</strong> {book.isbn}</p>
      {book.year && <p><strong>Year:</strong> {book.year}</p>}
      <p><strong>Publisher:</strong> {book.publisher?.text}</p>
      <p><strong>Genre:</strong> {book.genre?.text}</p>
      <p><strong>Condition:</strong> {book.condition?.text}</p>
      <p><strong>Type:</strong> {book.bookType?.text}</p>
      <p><strong>Stock:</strong> {book.quantity > 0 ? `${book.quantity} available` : 'Out of Stock'}</p>
      {book.summary && <p><strong>Summary:</strong> {book.summary}</p>}
      {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ maxWidth: '200px' }} />}
      <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem' }}>
        <button onClick={addToCart} disabled={book.quantity === 0}>Add to Cart</button>
        <button onClick={addToWishlist}>Add to Wishlist</button>
      </div>
    </div>
  );
}

export default BookDetailsPage;
