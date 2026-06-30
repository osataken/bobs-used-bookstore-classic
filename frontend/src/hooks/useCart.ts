import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getCart, deleteCartItem, getWishlist, moveToCart, moveAllToCart, deleteWishlistItem, addToCart, addToWishlist } from '../services/api';
import toast from 'react-hot-toast';

export function useCart() {
  return useQuery({
    queryKey: ['cart'],
    queryFn: async () => {
      const res = await getCart();
      return res.data;
    },
  });
}

export function useWishlist() {
  return useQuery({
    queryKey: ['wishlist'],
    queryFn: async () => {
      const res = await getWishlist();
      return res.data;
    },
  });
}

export function useAddToCart() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (bookId: number) => addToCart(bookId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      toast.success('Added to cart');
    },
  });
}

export function useAddToWishlist() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (bookId: number) => addToWishlist(bookId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      toast.success('Added to wishlist');
    },
  });
}

export function useDeleteCartItem() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => deleteCartItem(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      toast.success('Item removed');
    },
  });
}

export function useMoveToCart() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => moveToCart(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      toast.success('Moved to cart');
    },
  });
}

export function useMoveAllToCart() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => moveAllToCart(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['cart'] });
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      toast.success('All items moved to cart');
    },
  });
}

export function useDeleteWishlistItem() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => deleteWishlistItem(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['wishlist'] });
      toast.success('Item removed from wishlist');
    },
  });
}
