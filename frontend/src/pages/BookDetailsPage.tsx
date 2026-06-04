import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getBook, addToCart, addToWishlist } from '../services/api';

export default function BookDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const { data: book } = useQuery({ queryKey: ['book', id], queryFn: () => getBook(Number(id)) });

  if (!book) return <div>Loading...</div>;

  return (
    <div style={{ display: 'flex', gap: '2rem' }}>
      <div style={{ flex: '0 0 300px' }}>
        {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%' }} />}
      </div>
      <div>
        <h1>{book.name}</h1>
        <p><strong>Author:</strong> {book.author}</p>
        <p><strong>ISBN:</strong> {book.isbn}</p>
        {book.year && <p><strong>Year:</strong> {book.year}</p>}
        <p><strong>Publisher:</strong> {book.publisher?.text}</p>
        <p><strong>Genre:</strong> {book.genre?.text}</p>
        <p><strong>Condition:</strong> {book.condition?.text}</p>
        <p><strong>Type:</strong> {book.bookType?.text}</p>
        <p><strong>Price:</strong> ${book.price.toFixed(2)}</p>
        <p><strong>In Stock:</strong> {book.quantity > 0 ? `Yes (${book.quantity} available)` : 'No'}</p>
        {book.summary && <p>{book.summary}</p>}
        <div style={{ display: 'flex', gap: '1rem', marginTop: '1rem' }}>
          <button onClick={async () => { await addToCart(book.id); toast.success('Added to cart'); }} disabled={book.quantity === 0}>
            Add to Cart
          </button>
          <button onClick={async () => { await addToWishlist(book.id); toast.success('Added to wishlist'); }}>
            Add to Wishlist
          </button>
        </div>
      </div>
    </div>
  );
}
