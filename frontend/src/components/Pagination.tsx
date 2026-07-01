interface PaginationProps {
  pageIndex: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

function Pagination({ pageIndex, totalPages, onPageChange }: PaginationProps) {
  if (totalPages <= 1) return null;

  return (
    <div style={{ display: 'flex', gap: '0.5rem', justifyContent: 'center', margin: '1rem 0' }}>
      <button disabled={pageIndex <= 1} onClick={() => onPageChange(pageIndex - 1)}>Previous</button>
      <span>Page {pageIndex} of {totalPages}</span>
      <button disabled={pageIndex >= totalPages} onClick={() => onPageChange(pageIndex + 1)}>Next</button>
    </div>
  );
}

export default Pagination;
