import { describe, it, expect } from 'vitest';
import { OrderStatus, OfferStatus, ReferenceDataType } from './index';

describe('Types', () => {
  it('OrderStatus enum values are correct', () => {
    expect(OrderStatus.Pending).toBe(0);
    expect(OrderStatus.Ordered).toBe(1);
    expect(OrderStatus.Shipped).toBe(2);
    expect(OrderStatus.Delivered).toBe(3);
    expect(OrderStatus.Cancelled).toBe(4);
  });

  it('OfferStatus enum values are correct', () => {
    expect(OfferStatus.PendingApproval).toBe(0);
    expect(OfferStatus.Approved).toBe(1);
    expect(OfferStatus.Received).toBe(2);
    expect(OfferStatus.Paid).toBe(3);
    expect(OfferStatus.Rejected).toBe(4);
  });

  it('ReferenceDataType enum values are correct', () => {
    expect(ReferenceDataType.Publisher).toBe(0);
    expect(ReferenceDataType.Condition).toBe(1);
    expect(ReferenceDataType.BookType).toBe(2);
    expect(ReferenceDataType.Genre).toBe(3);
  });
});
