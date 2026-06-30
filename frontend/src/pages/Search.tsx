import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useSearchBooks } from '../hooks/useBooks'
import { useAddToCart, useAddToWishlist } from '../hooks/useCart'

function Search() {
  const [searchString, setSearchString] = useState('')
  const [sortBy, setSortBy] = useState('Name')
  const [pageIndex, setPageIndex] = useState(1)
  const pageSize = 10

  const { data, isLoading } = useSearchBooks({ searchString, sortBy, pageIndex, pageSize })
  const addToCart = useAddToCart()
  const addToWishlist = useAddToWishlist()

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
        <select value={sortBy} onChange={e => setSortBy(e.target.value)}>
          <option value="Name">Sort by Name</option>
          <option value="Author">Sort by Author</option>
          <option value="Price">Sort by Price</option>
        </select>
      </div>

      {isLoading && <p>Loading...</p>}

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '1rem' }}>
        {data?.items.map(book => (
          <div key={book.id} style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
            {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%', height: '150px', objectFit: 'cover' }} />}
            <h3><Link to={`/books/${book.id}`}>{book.name}</Link></h3>
            <p>{book.author}</p>
            <p>${book.price.toFixed(2)} {book.quantity <= 5 && book.quantity > 0 && <span style={{ color: 'orange' }}>Low Stock</span>}</p>
            <div style={{ display: 'flex', gap: '0.5rem' }}>
              <button onClick={() => addToCart.mutate(book.id)} disabled={book.quantity === 0}>
                {book.quantity > 0 ? 'Add to Cart' : 'Out of Stock'}
              </button>
              <button onClick={() => addToWishlist.mutate(book.id)}>Wishlist</button>
            </div>
          </div>
        ))}
      </div>

      {data && data.totalPages > 1 && (
        <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
          <button onClick={() => setPageIndex(p => Math.max(1, p - 1))} disabled={pageIndex === 1}>Previous</button>
          <span>Page {pageIndex} of {data.totalPages}</span>
          <button onClick={() => setPageIndex(p => Math.min(data.totalPages, p + 1))} disabled={pageIndex === data.totalPages}>Next</button>
        </div>
      )}
    </div>
  )
}

export default Search
