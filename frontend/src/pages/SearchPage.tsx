import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { searchBooks, addToCart, addToWishlist } from '../services/api'
import Pagination from '../components/Pagination'
import toast from 'react-hot-toast'

export default function SearchPage() {
  const [searchString, setSearchString] = useState('')
  const [sortBy, setSortBy] = useState('Name')
  const [pageIndex, setPageIndex] = useState(1)

  const { data, isLoading } = useQuery({
    queryKey: ['search', searchString, sortBy, pageIndex],
    queryFn: () => searchBooks({ searchString, sortBy, pageIndex, pageSize: 10 }).then(r => r.data),
  })

  const handleAddToCart = async (bookId: number) => {
    try {
      await addToCart(bookId)
      toast.success('Added to cart')
    } catch {
      toast.error('Failed to add to cart')
    }
  }

  const handleAddToWishlist = async (bookId: number) => {
    try {
      await addToWishlist(bookId)
      toast.success('Added to wishlist')
    } catch {
      toast.error('Failed to add to wishlist')
    }
  }

  return (
    <div>
      <h1>Search Books</h1>
      <div style={{ display: 'flex', gap: '10px', marginBottom: '20px' }}>
        <input
          type="text"
          placeholder="Search..."
          value={searchString}
          onChange={e => { setSearchString(e.target.value); setPageIndex(1) }}
          style={{ padding: '8px', flex: 1 }}
        />
        <select value={sortBy} onChange={e => setSortBy(e.target.value)} style={{ padding: '8px' }}>
          <option value="Name">Name</option>
          <option value="PriceAsc">Price (Low to High)</option>
          <option value="PriceDesc">Price (High to Low)</option>
        </select>
      </div>

      {isLoading ? <div>Loading...</div> : (
        <>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '20px' }}>
            {data?.items?.map(book => (
              <div key={book.id} style={{ border: '1px solid #ddd', padding: '15px', borderRadius: '8px' }}>
                {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%', height: '200px', objectFit: 'cover' }} />}
                <h3><Link to={`/search/${book.id}`}>{book.name}</Link></h3>
                <p>{book.author}</p>
                <p>{book.genre} | {book.condition}</p>
                <p><strong>${book.price.toFixed(2)}</strong></p>
                <p>{book.isInStock ? (book.isLowInStock ? 'Low Stock' : 'In Stock') : 'Out of Stock'}</p>
                <div style={{ display: 'flex', gap: '5px' }}>
                  <button onClick={() => handleAddToCart(book.id)} disabled={!book.isInStock}>Add to Cart</button>
                  <button onClick={() => handleAddToWishlist(book.id)}>Wishlist</button>
                </div>
              </div>
            ))}
          </div>
          {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
        </>
      )}
    </div>
  )
}
