import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import api from '../services/api';
import { Book } from '../types';

function HomePage() {
  const { data, isLoading } = useQuery({
    queryKey: ['bestSellers'],
    queryFn: async () => {
      const res = await api.get('/home');
      return res.data.bestSellers as Book[];
    },
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Welcome to Bob's Used Bookstore</h1>
      <h2>Best Sellers</h2>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '1rem' }}>
        {data?.map((book) => (
          <div key={book.id} style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
            <Link to={`/search/${book.id}`}>
              <h3>{book.name}</h3>
            </Link>
            <p>by {book.author}</p>
            <p>${book.price.toFixed(2)}</p>
            <p>{book.quantity > 0 ? 'In Stock' : 'Out of Stock'}</p>
          </div>
        ))}
      </div>
    </div>
  );
}

export default HomePage;
