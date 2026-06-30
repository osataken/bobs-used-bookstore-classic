import { useParams } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { getBookDetails, addToCart, addToWishlist } from '../services/api'
import toast from 'react-hot-toast'

export default function BookDetailsPage() {
  const { id } = useParams<{ id: string }>()
  const { data: book, isLoading } = useQuery({
    queryKey: ['book', id],
    queryFn: () => getBookDetails(Number(id)).then(r => r.data),
    enabled: !!id,
  })

  if (isLoading) return <div>Loading...</div>
  if (!book) return <div>Book not found</div>

  const handleAddToCart = async () => {
    try {
      await addToCart(book.id)
      toast.success('Added to cart')
    } catch {
      toast.error('Failed to add to cart')
    }
  }

  const handleAddToWishlist = async () => {
    try {
      await addToWishlist(book.id)
      toast.success('Added to wishlist')
    } catch {
      toast.error('Failed to add to wishlist')
    }
  }

  return (
    <div>
      <h1>{book.name}</h1>
      <div style={{ display: 'flex', gap: '30px' }}>
        {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '300px', height: '450px', objectFit: 'cover' }} />}
        <div>
          <p><strong>Author:</strong> {book.author}</p>
          <p><strong>ISBN:</strong> {book.isbn}</p>
          {book.year && <p><strong>Year:</strong> {book.year}</p>}
          <p><strong>Genre:</strong> {book.genre}</p>
          <p><strong>Type:</strong> {book.bookType}</p>
          <p><strong>Publisher:</strong> {book.publisher}</p>
          <p><strong>Condition:</strong> {book.condition}</p>
          <p><strong>Price:</strong> ${book.price.toFixed(2)}</p>
          <p><strong>Status:</strong> {book.isInStock ? (book.isLowInStock ? 'Low Stock' : 'In Stock') : 'Out of Stock'}</p>
          {book.summary && <p><strong>Summary:</strong> {book.summary}</p>}
          <div style={{ display: 'flex', gap: '10px', marginTop: '15px' }}>
            <button onClick={handleAddToCart} disabled={!book.isInStock}>Add to Cart</button>
            <button onClick={handleAddToWishlist}>Add to Wishlist</button>
          </div>
        </div>
      </div>
    </div>
  )
}
