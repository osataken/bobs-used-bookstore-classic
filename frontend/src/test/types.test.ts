import { describe, it, expect } from 'vitest';
import { OrderStatusLabels, OfferStatusLabels, ReferenceDataTypeLabels } from '../types';

describe('Types', () => {
  it('OrderStatusLabels has correct values', () => {
    expect(OrderStatusLabels[0]).toBe('Pending');
    expect(OrderStatusLabels[1]).toBe('Ordered');
    expect(OrderStatusLabels[2]).toBe('Shipped');
    expect(OrderStatusLabels[3]).toBe('Delivered');
    expect(OrderStatusLabels[4]).toBe('Cancelled');
  });

  it('OfferStatusLabels has correct values', () => {
    expect(OfferStatusLabels[0]).toBe('Pending Approval');
    expect(OfferStatusLabels[1]).toBe('Approved');
    expect(OfferStatusLabels[2]).toBe('Received');
    expect(OfferStatusLabels[3]).toBe('Paid');
    expect(OfferStatusLabels[4]).toBe('Rejected');
  });

  it('ReferenceDataTypeLabels has correct values', () => {
    expect(ReferenceDataTypeLabels[0]).toBe('Publisher');
    expect(ReferenceDataTypeLabels[1]).toBe('Condition');
    expect(ReferenceDataTypeLabels[2]).toBe('BookType');
    expect(ReferenceDataTypeLabels[3]).toBe('Genre');
  });
});
