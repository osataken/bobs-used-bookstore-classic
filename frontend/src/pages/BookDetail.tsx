import { useParams } from 'react-router-dom'
import { useBookDetails } from '../hooks/useBooks'
import { useAddToCart, useAddToWishlist } from '../hooks/useCart'

function BookDetail() {
  const { id } = useParams<{ id: string }>()
  const { data: book, isLoading } = useBookDetails(Number(id))
  const addToCart = useAddToCart()
  const addToWishlist = useAddToWishlist()

  if (isLoading) return <p>Loading...</p>
  if (!book) return <p>Book not found</p>

  return (
    <div style={{ display: 'flex', gap: '2rem' }}>
      <div>
        {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '300px', borderRadius: '8px' }} />}
      </div>
      <div>
        <h1>{book.name}</h1>
        <p><strong>Author:</strong> {book.author}</p>
        <p><strong>ISBN:</strong> {book.isbn}</p>
        {book.year && <p><strong>Year:</strong> {book.year}</p>}
        <p><strong>Price:</strong> ${book.price.toFixed(2)}</p>
        <p><strong>Condition:</strong> {book.condition?.text}</p>
        <p><strong>Genre:</strong> {book.genre?.text}</p>
        <p><strong>Publisher:</strong> {book.publisher?.text}</p>
        <p><strong>Type:</strong> {book.bookType?.text}</p>
        <p><strong>In Stock:</strong> {book.quantity > 0 ? `Yes (${book.quantity} available)` : 'No'}</p>
        {book.summary && <p><strong>Summary:</strong> {book.summary}</p>}
        <div style={{ display: 'flex', gap: '0.5rem', marginTop: '1rem' }}>
          <button onClick={() => addToCart.mutate(book.id)} disabled={book.quantity === 0}>
            {book.quantity > 0 ? 'Add to Cart' : 'Out of Stock'}
          </button>
          <button onClick={() => addToWishlist.mutate(book.id)}>Add to Wishlist</button>
        </div>
      </div>
    </div>
  )
}

export default BookDetail
