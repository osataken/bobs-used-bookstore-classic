import { useQuery } from '@tanstack/react-query';
import { searchBooks, getBookDetails, getHome } from '../services/api';

export function useBestSellers() {
  return useQuery({
    queryKey: ['bestSellers'],
    queryFn: async () => {
      const res = await getHome();
      return res.data.bestSellers;
    },
  });
}

export function useSearchBooks(params: { searchString?: string; sortBy?: string; pageIndex?: number; pageSize?: number }) {
  return useQuery({
    queryKey: ['books', params],
    queryFn: async () => {
      const res = await searchBooks(params);
      return res.data;
    },
  });
}

export function useBookDetails(id: number) {
  return useQuery({
    queryKey: ['book', id],
    queryFn: async () => {
      const res = await getBookDetails(id);
      return res.data;
    },
    enabled: id > 0,
  });
}
