using Bookstore.Domain.Customers;

namespace Bookstore.Domain.Addresses
{
    public class Address : Entity
    {
        // An empty constructor is required by EF Core
        private Address() { }

        public Address(Customer customer, string addressLine1, string addressLine2, string city, string state, string country, string zipCode)
        {
            Customer = customer;
            CustomerId = customer.Id;
            AddressLine1 = addressLine1;
            AddressLine2 = addressLine2;
            City = city;
            State = state;
            Country = country;
            ZipCode = zipCode;
        }

        public string AddressLine1 { get; set; } = string.Empty;

        public string? AddressLine2 { get; set; }

        public string City { get; set; } = string.Empty;

        public string State { get; set; } = string.Empty;

        public string Country { get; set; } = string.Empty;

        public string ZipCode { get; set; } = string.Empty;

        public int CustomerId { get; set; }
        public Customer Customer { get; set; } = null!;

        public bool IsActive { get; set; } = true;
    }
}