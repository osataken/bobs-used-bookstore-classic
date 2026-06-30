import { useState } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { getAdminOffers, approveOffer, rejectOffer, receivedOffer, paidOffer } from '../../services/api'
import Pagination from '../../components/Pagination'
import toast from 'react-hot-toast'

export default function OffersPage() {
  const queryClient = useQueryClient()
  const [pageIndex, setPageIndex] = useState(1)

  const { data, isLoading } = useQuery({
    queryKey: ['adminOffers', pageIndex],
    queryFn: () => getAdminOffers({ pageIndex, pageSize: 10 }).then(r => r.data),
  })

  const handleAction = async (action: (id: number) => Promise<unknown>, id: number, label: string) => {
    try {
      await action(id)
      queryClient.invalidateQueries({ queryKey: ['adminOffers'] })
      toast.success(`Offer ${label}`)
    } catch {
      toast.error(`Failed to ${label} offer`)
    }
  }

  if (isLoading) return <div>Loading...</div>

  return (
    <div>
      <h1>Offer Management</h1>
      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        <thead>
          <tr>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Book</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Author</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Price</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Status</th>
            <th style={{ textAlign: 'left', padding: '8px', borderBottom: '1px solid #ddd' }}>Actions</th>
          </tr>
        </thead>
        <tbody>
          {data?.items?.map(offer => (
            <tr key={offer.id}>
              <td style={{ padding: '8px' }}>{offer.bookName}</td>
              <td style={{ padding: '8px' }}>{offer.author}</td>
              <td style={{ padding: '8px' }}>${offer.bookPrice.toFixed(2)}</td>
              <td style={{ padding: '8px' }}>{offer.statusText}</td>
              <td style={{ padding: '8px' }}>
                {offer.offerStatus === 0 && (
                  <>
                    <button onClick={() => handleAction(approveOffer, offer.id, 'approved')}>Approve</button>{' '}
                    <button onClick={() => handleAction(rejectOffer, offer.id, 'rejected')}>Reject</button>
                  </>
                )}
                {offer.offerStatus === 1 && <button onClick={() => handleAction(receivedOffer, offer.id, 'received')}>Mark Received</button>}
                {offer.offerStatus === 2 && <button onClick={() => handleAction(paidOffer, offer.id, 'paid')}>Mark Paid</button>}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      {data && <Pagination pageIndex={data.pageIndex} totalPages={data.totalPages} onPageChange={setPageIndex} />}
    </div>
  )
}
