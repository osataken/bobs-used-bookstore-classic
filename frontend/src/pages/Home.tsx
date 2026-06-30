import { Link } from 'react-router-dom'
import { useBestSellers } from '../hooks/useBooks'
import { useAddToCart } from '../hooks/useCart'

function Home() {
  const { data: books, isLoading } = useBestSellers()
  const addToCart = useAddToCart()

  if (isLoading) return <p>Loading...</p>

  return (
    <div>
      <h1>Welcome to Bob's Used Bookstore</h1>
      <h2>Best Sellers</h2>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '1rem' }}>
        {books?.map(book => (
          <div key={book.id} style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
            {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%', height: '200px', objectFit: 'cover' }} />}
            <h3><Link to={`/books/${book.id}`}>{book.name}</Link></h3>
            <p>{book.author}</p>
            <p style={{ fontWeight: 'bold' }}>${book.price.toFixed(2)}</p>
            <button onClick={() => addToCart.mutate(book.id)} disabled={book.quantity === 0}>
              {book.quantity > 0 ? 'Add to Cart' : 'Out of Stock'}
            </button>
          </div>
        ))}
      </div>
    </div>
  )
}

export default Home
