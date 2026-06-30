interface PaginationProps {
  pageIndex: number
  totalPages: number
  onPageChange: (page: number) => void
}

export default function Pagination({ pageIndex, totalPages, onPageChange }: PaginationProps) {
  if (totalPages <= 1) return null

  const pages = Array.from({ length: totalPages }, (_, i) => i + 1)

  return (
    <div style={{ display: 'flex', gap: '5px', marginTop: '20px', justifyContent: 'center' }}>
      <button disabled={pageIndex <= 1} onClick={() => onPageChange(pageIndex - 1)}>Previous</button>
      {pages.map(p => (
        <button key={p} onClick={() => onPageChange(p)} style={{ fontWeight: p === pageIndex ? 'bold' : 'normal' }}>
          {p}
        </button>
      ))}
      <button disabled={pageIndex >= totalPages} onClick={() => onPageChange(pageIndex + 1)}>Next</button>
    </div>
  )
}
