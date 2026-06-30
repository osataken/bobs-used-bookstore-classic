import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getOrders, getOrderDetails, cancelOrder } from '../services/api';
import toast from 'react-hot-toast';

export function useOrders() {
  return useQuery({
    queryKey: ['orders'],
    queryFn: async () => {
      const res = await getOrders();
      return res.data.orders;
    },
  });
}

export function useOrderDetails(id: number) {
  return useQuery({
    queryKey: ['order', id],
    queryFn: async () => {
      const res = await getOrderDetails(id);
      return res.data;
    },
    enabled: id > 0,
  });
}

export function useCancelOrder() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => cancelOrder(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      toast.success('Order cancelled');
    },
  });
}
