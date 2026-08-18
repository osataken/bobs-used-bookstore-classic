import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { homeService } from '../services';
import { Book } from '../types';

export default function Home() {
  const { data: books, isLoading } = useQuery({
    queryKey: ['featured'],
    queryFn: () => homeService.getFeatured().then(r => r.data),
  });

  return (
    <div>
      <h1 className="page-title">Bob's Used Bookstore</h1>
      <p>Welcome! Browse our collection of used books.</p>
      <h2 style={{ margin: '1rem 0' }}>Featured Books</h2>
      {isLoading ? <p>Loading...</p> : (
        <div className="grid">
          {books?.map((book: Book) => (
            <div key={book.id} className="card">
              {book.coverImageUrl && <img src={book.coverImageUrl} alt={book.name} style={{ width: '100%', height: 200, objectFit: 'cover' }} />}
              <h3>{book.name}</h3>
              <p>by {book.author}</p>
              <p><strong>${book.price.toFixed(2)}</strong></p>
              <Link to={`/search/${book.id}`}>View Details</Link>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
