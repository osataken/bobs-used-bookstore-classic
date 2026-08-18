import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { searchService, cartService, wishlistService } from '../services';
import toast from 'react-hot-toast';

export default function BookDetails() {
  const { id } = useParams<{ id: string }>();
  const { data: book, isLoading } = useQuery({
    queryKey: ['book', id],
    queryFn: () => searchService.getBook(Number(id)).then(r => r.data),
  });

  if (isLoading) return <p>Loading...</p>;
  if (!book) return <p>Book not found</p>;

  return (
    <div>
      <h1 className="page-title">{book.name}</h1>
      <div style={{ display: 'flex', gap: '2rem' }}>
        {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: 200, height: 300, objectFit: 'cover' }} />}
        <div>
          <p><strong>Author:</strong> {book.author}</p>
          <p><strong>ISBN:</strong> {book.isbn}</p>
          <p><strong>Price:</strong> ${book.price.toFixed(2)}</p>
          <p><strong>Stock:</strong> {book.quantity > 0 ? `${book.quantity} available` : 'Out of stock'}</p>
          {book.publisher && <p><strong>Publisher:</strong> {book.publisher.text}</p>}
          {book.genre && <p><strong>Genre:</strong> {book.genre.text}</p>}
          {book.condition && <p><strong>Condition:</strong> {book.condition.text}</p>}
          {book.bookType && <p><strong>Type:</strong> {book.bookType.text}</p>}
          {book.summary && <p><strong>Summary:</strong> {book.summary}</p>}
          <div style={{ marginTop: '1rem', display: 'flex', gap: '8px' }}>
            <button className="btn btn-primary" disabled={book.quantity === 0} onClick={async () => { await cartService.addItem(book.id); toast.success('Added to cart'); }}>Add to Cart</button>
            <button className="btn" onClick={async () => { await wishlistService.addItem(book.id); toast.success('Added to wishlist'); }}>Add to Wishlist</button>
          </div>
        </div>
      </div>
    </div>
  );
}
