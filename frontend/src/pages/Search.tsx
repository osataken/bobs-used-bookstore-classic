import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { searchService, cartService, wishlistService } from '../services';
import { Book } from '../types';
import toast from 'react-hot-toast';

export default function Search() {
  const [searchString, setSearchString] = useState('');
  const [sortBy, setSortBy] = useState('');
  const [pageIndex, setPageIndex] = useState(1);
  const pageSize = 10;

  const { data, isLoading } = useQuery({
    queryKey: ['search', searchString, sortBy, pageIndex],
    queryFn: () => searchService.search({ searchString, sortBy, pageIndex, pageSize }).then(r => r.data),
  });

  const addToCart = async (bookId: number) => {
    try {
      await cartService.addItem(bookId);
      toast.success('Added to cart');
    } catch { toast.error('Failed to add to cart'); }
  };

  const addToWishlist = async (bookId: number) => {
    try {
      await wishlistService.addItem(bookId);
      toast.success('Added to wishlist');
    } catch { toast.error('Failed to add to wishlist'); }
  };

  return (
    <div>
      <h1 className="page-title">Browse Books</h1>
      <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
        <input placeholder="Search by title, author, or ISBN..." value={searchString} onChange={e => { setSearchString(e.target.value); setPageIndex(1); }} />
        <select value={sortBy} onChange={e => setSortBy(e.target.value)}>
          <option value="">Sort by</option>
          <option value="price_asc">Price: Low to High</option>
          <option value="price_desc">Price: High to Low</option>
          <option value="name_asc">Name: A-Z</option>
          <option value="name_desc">Name: Z-A</option>
        </select>
      </div>
      {isLoading ? <p>Loading...</p> : (
        <>
          <div className="grid">
            {data?.items?.map((book: Book) => (
              <div key={book.id} className="card">
                {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%', height: 150, objectFit: 'cover' }} />}
                <h3><Link to={`/search/${book.id}`}>{book.name}</Link></h3>
                <p>by {book.author}</p>
                <p>${book.price.toFixed(2)} {book.quantity > 0 ? <span style={{ color: 'green' }}>In Stock</span> : <span style={{ color: 'red' }}>Out of Stock</span>}</p>
                <div style={{ display: 'flex', gap: '4px', marginTop: '8px' }}>
                  <button className="btn btn-primary" onClick={() => addToCart(book.id)} disabled={book.quantity === 0}>Add to Cart</button>
                  <button className="btn" onClick={() => addToWishlist(book.id)}>Wishlist</button>
                </div>
              </div>
            ))}
          </div>
          {data && data.totalPages > 1 && (
            <div className="pagination">
              <button className="btn" disabled={pageIndex <= 1} onClick={() => setPageIndex(p => p - 1)}>Previous</button>
              <span>Page {data.pageIndex} of {data.totalPages}</span>
              <button className="btn" disabled={pageIndex >= data.totalPages} onClick={() => setPageIndex(p => p + 1)}>Next</button>
            </div>
          )}
        </>
      )}
    </div>
  );
}
