import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { searchBooks, addToCart, addToWishlist } from '../services/api';

export default function SearchPage() {
  const [searchString, setSearchString] = useState('');
  const [sortBy, setSortBy] = useState('Name');
  const [pageIndex, setPageIndex] = useState(1);
  const pageSize = 10;

  const { data } = useQuery({
    queryKey: ['books', searchString, sortBy, pageIndex],
    queryFn: () => searchBooks({ searchString, sortBy, pageIndex, pageSize }),
  });

  const handleAddToCart = async (bookId: number) => {
    await addToCart(bookId);
    toast.success('Item added to shopping cart');
  };

  const handleAddToWishlist = async (bookId: number) => {
    await addToWishlist(bookId);
    toast.success('Item added to wishlist');
  };

  return (
    <div>
      <h1>Browse Books</h1>
      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
        <input
          type="text"
          placeholder="Search by name, author, or ISBN..."
          value={searchString}
          onChange={e => { setSearchString(e.target.value); setPageIndex(1); }}
          style={{ flex: 1, padding: '0.5rem' }}
        />
        <select value={sortBy} onChange={e => setSortBy(e.target.value)} style={{ padding: '0.5rem' }}>
          <option value="Name">Sort by Name</option>
          <option value="Price">Sort by Price</option>
          <option value="Author">Sort by Author</option>
        </select>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '1rem' }}>
        {(data?.items as Array<{ id: number; name: string; author: string; price: number; quantity: number; coverImageUrl: string }>)?.map(book => (
          <div key={book.id} style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
            {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%', height: '150px', objectFit: 'cover' }} />}
            <h3><Link to={`/books/${book.id}`}>{book.name}</Link></h3>
            <p>{book.author}</p>
            <p><strong>${book.price.toFixed(2)}</strong> {book.quantity === 0 && <span style={{ color: 'red' }}>(Out of Stock)</span>}</p>
            <div style={{ display: 'flex', gap: '0.5rem' }}>
              <button onClick={() => handleAddToCart(book.id)} disabled={book.quantity === 0}>Add to Cart</button>
              <button onClick={() => handleAddToWishlist(book.id)}>Wishlist</button>
            </div>
          </div>
        ))}
      </div>

      {data && data.totalPages > 1 && (
        <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
          <button disabled={pageIndex <= 1} onClick={() => setPageIndex(p => p - 1)}>Previous</button>
          <span>Page {pageIndex} of {data.totalPages}</span>
          <button disabled={pageIndex >= data.totalPages} onClick={() => setPageIndex(p => p + 1)}>Next</button>
        </div>
      )}
    </div>
  );
}
