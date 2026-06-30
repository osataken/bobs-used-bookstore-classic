import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { getAdminOffers, approveOffer, rejectOffer, markOfferReceived, markOfferPaid } from '../../services/api'
import toast from 'react-hot-toast'

function AdminOffers() {
  const queryClient = useQueryClient()
  const [pageIndex, setPageIndex] = useState(1)

  const { data } = useQuery({
    queryKey: ['adminOffers', pageIndex],
    queryFn: async () => { const res = await getAdminOffers({ pageIndex, pageSize: 10 }); return res.data; },
  })

  const approve = useMutation({ mutationFn: (id: number) => approveOffer(id), onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminOffers'] }); toast.success('Offer approved'); } })
  const reject = useMutation({ mutationFn: (id: number) => rejectOffer(id), onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminOffers'] }); toast.success('Offer rejected'); } })
  const received = useMutation({ mutationFn: (id: number) => markOfferReceived(id), onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminOffers'] }); toast.success('Marked received'); } })
  const paid = useMutation({ mutationFn: (id: number) => markOfferPaid(id), onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['adminOffers'] }); toast.success('Marked paid'); } })

  const statusLabels = ['Pending Approval', 'Approved', 'Received', 'Paid', 'Rejected']

  return (
    <div>
      <h1>Offer Management</h1>
      {data && (
        <>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead><tr><th style={{ textAlign: 'left' }}>Book</th><th style={{ textAlign: 'left' }}>Author</th><th style={{ textAlign: 'left' }}>Customer</th><th style={{ textAlign: 'left' }}>Price</th><th style={{ textAlign: 'left' }}>Status</th><th>Actions</th></tr></thead>
            <tbody>
              {data.items.map(offer => (
                <tr key={offer.id} style={{ borderBottom: '1px solid #ddd' }}>
                  <td style={{ padding: '0.5rem' }}>{offer.bookName}</td>
                  <td style={{ padding: '0.5rem' }}>{offer.author}</td>
                  <td style={{ padding: '0.5rem' }}>{offer.customerId}</td>
                  <td style={{ padding: '0.5rem' }}>${offer.bookPrice.toFixed(2)}</td>
                  <td style={{ padding: '0.5rem' }}>{statusLabels[offer.offerStatus]}</td>
                  <td style={{ padding: '0.5rem', display: 'flex', gap: '0.25rem', flexWrap: 'wrap' }}>
                    {offer.offerStatus === 0 && <><button onClick={() => approve.mutate(offer.id)}>Approve</button><button onClick={() => reject.mutate(offer.id)}>Reject</button></>}
                    {offer.offerStatus === 1 && <button onClick={() => received.mutate(offer.id)}>Received</button>}
                    {offer.offerStatus === 2 && <button onClick={() => paid.mutate(offer.id)}>Paid</button>}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          {data.totalPages > 1 && (
            <div style={{ marginTop: '1rem', display: 'flex', gap: '0.5rem', justifyContent: 'center' }}>
              <button onClick={() => setPageIndex(p => Math.max(1, p - 1))} disabled={pageIndex === 1}>Previous</button>
              <span>Page {pageIndex} of {data.totalPages}</span>
              <button onClick={() => setPageIndex(p => p + 1)} disabled={pageIndex === data.totalPages}>Next</button>
            </div>
          )}
        </>
      )}
    </div>
  )
}

export default AdminOffers
