import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getAdminInventory } from '../../services/api'
import Pagination from '../../components/Pagination'
import { Link } from 'react-router-dom'

export default function InventoryPage() {
  const [pageIndex, setPageIndex] = useState(1)
  const [name, setName] = useState('')

  const { data, isLoading } = useQuery({
    queryKey: ['adminInventory', pageIndex, name],
    queryFn: () => getAdminInventory({ pageIndex, pageSize: 10, name }).then(r => r.data),
  })

  return (
    <div>
      <h1>Inventory Management</h1>
      <div style={{ display: 'flex', gap: '10px', marginBottom: '15px' }}>
        <input placeholder="Search by name" value={name} onChange={e => { setName(e.target.value); setPageIndex(1) }} />
      </div>
      {isLoading ? <div>Loading...</div> : (
        <>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Name</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Author</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Price</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Qty</th>
                <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {data?.items?.map(book => (
                <tr key={book.id}>
                  <td style={{ padding: '8px' }}>{book.name}</td>
                  <td style={{ padding: '8px' }}>{book.author}</td>
                  <td style={{ padding: '8px' }}>${book.price.toFixed(2)}</td>
                  <td style={{ padding: '8px' }}>{book.quantity}</td>
                  <td style={{ padding: '8px' }}><Link to={`/admin/inventory/${book.id}`}>View</Link></td>
                </tr>
              ))}
            </tbody>
          </table>
          {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
        </>
      )}
    </div>
  )
}
