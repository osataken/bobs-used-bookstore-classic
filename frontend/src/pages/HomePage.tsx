import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { getBestSellingBooks } from '../services/api';

export default function HomePage() {
  const { data: books } = useQuery({ queryKey: ['best-selling'], queryFn: () => getBestSellingBooks(4) });

  return (
    <div>
      <h1>Welcome to Bob's Used Bookstore</h1>
      <p>Find your next great read at a great price.</p>
      <h2>Popular Books</h2>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '1rem' }}>
        {books?.map(book => (
          <div key={book.id} style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
            {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%', height: '200px', objectFit: 'cover' }} />}
            <h3><Link to={`/books/${book.id}`}>{book.name}</Link></h3>
            <p>{book.author}</p>
            <p><strong>${book.price.toFixed(2)}</strong></p>
          </div>
        ))}
      </div>
      <p style={{ marginTop: '2rem' }}>
        <Link to="/search">Browse all books →</Link>
      </p>
    </div>
  );
}
