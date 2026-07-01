import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import api from '../services/api';
import { Book, PaginatedResponse } from '../types';
import Pagination from '../components/Pagination';

function SearchPage() {
  const [searchString, setSearchString] = useState('');
  const [sortBy, setSortBy] = useState('Name');
  const [pageIndex, setPageIndex] = useState(1);
  const pageSize = 10;

  const { data, isLoading } = useQuery({
    queryKey: ['search', searchString, sortBy, pageIndex],
    queryFn: async () => {
      const res = await api.get('/search', {
        params: { searchString, sortBy, pageIndex, pageSize },
      });
      return res.data as PaginatedResponse<Book>;
    },
  });

  const addToCart = async (bookId: number) => {
    try {
      await api.post('/cart/add', { bookId, quantity: 1 });
      toast.success('Added to cart');
    } catch {
      toast.error('Failed to add to cart');
    }
  };

  const addToWishlist = async (bookId: number) => {
    try {
      await api.post('/wishlist/add', { bookId });
      toast.success('Added to wishlist');
    } catch {
      toast.error('Failed to add to wishlist');
    }
  };

  return (
    <div>
      <h1>Browse Books</h1>
      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
        <input
          type="text"
          placeholder="Search by name, author, or ISBN..."
          value={searchString}
          onChange={(e) => { setSearchString(e.target.value); setPageIndex(1); }}
          style={{ flex: 1, padding: '0.5rem' }}
        />
        <select value={sortBy} onChange={(e) => setSortBy(e.target.value)}>
          <option value="Name">Sort by Name</option>
          <option value="Author">Sort by Author</option>
          <option value="Price">Sort by Price</option>
        </select>
      </div>

      {isLoading ? (
        <div>Loading...</div>
      ) : (
        <>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '1rem' }}>
            {data?.items.map((book) => (
              <div key={book.id} style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
                <Link to={`/search/${book.id}`}><h3>{book.name}</h3></Link>
                <p>by {book.author}</p>
                <p>${book.price.toFixed(2)}</p>
                <p>{book.quantity > 0 ? `In Stock (${book.quantity})` : 'Out of Stock'}</p>
                <div style={{ display: 'flex', gap: '0.5rem' }}>
                  <button onClick={() => addToCart(book.id)} disabled={book.quantity === 0}>Add to Cart</button>
                  <button onClick={() => addToWishlist(book.id)}>Wishlist</button>
                </div>
              </div>
            ))}
          </div>
          {data && (
            <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />
          )}
        </>
      )}
    </div>
  );
}

export default SearchPage;
