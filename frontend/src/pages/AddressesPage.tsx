import { useQuery, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { getAddresses, deleteAddress } from '../services/api';

export default function AddressesPage() {
  const queryClient = useQueryClient();
  const { data: addresses } = useQuery({ queryKey: ['addresses'], queryFn: getAddresses });

  const handleDelete = async (id: number) => {
    await deleteAddress(id);
    toast.success('Address deleted');
    queryClient.invalidateQueries({ queryKey: ['addresses'] });
  };

  return (
    <div>
      <h1>My Addresses</h1>
      <Link to="/addresses/new"><button style={{ marginBottom: '1rem' }}>Add New Address</button></Link>
      {!addresses?.length ? (
        <p>No addresses yet.</p>
      ) : (
        <div style={{ display: 'grid', gap: '1rem' }}>
          {addresses.map(addr => (
            <div key={addr.id} style={{ border: '1px solid #ddd', padding: '1rem', borderRadius: '8px' }}>
              <p>{addr.addressLine1}</p>
              {addr.addressLine2 && <p>{addr.addressLine2}</p>}
              <p>{addr.city}, {addr.state} {addr.zipCode}</p>
              <p>{addr.country}</p>
              <div style={{ display: 'flex', gap: '0.5rem' }}>
                <Link to={`/addresses/${addr.id}/edit`}><button>Edit</button></Link>
                <button onClick={() => handleDelete(addr.id)} style={{ color: 'red' }}>Delete</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
