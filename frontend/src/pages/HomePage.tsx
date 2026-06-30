import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getHome } from '../services/api'

export default function HomePage() {
  const { data, isLoading } = useQuery({
    queryKey: ['home'],
    queryFn: () => getHome().then(r => r.data),
  })

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      <h1>Welcome to Bob's Used Bookstore</h1>
      <h2>Best Selling Books</h2>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '20px' }}>
        {data?.bestSellingBooks?.map(book => (
          <div key={book.id} style={{ border: '1px solid #ddd', padding: '15px', borderRadius: '8px' }}>
            {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%', height: '200px', objectFit: 'cover' }} />}
            <h3><Link to={`/search/${book.id}`}>{book.name}</Link></h3>
            <p>{book.author}</p>
            <p><strong>${book.price.toFixed(2)}</strong></p>
            <p>{book.isInStock ? 'In Stock' : 'Out of Stock'}</p>
          </div>
        ))}
      </div>
      <div style={{ marginTop: '20px' }}>
        <Link to="/search">Browse All Books</Link>
      </div>
    </div>
  )
}
