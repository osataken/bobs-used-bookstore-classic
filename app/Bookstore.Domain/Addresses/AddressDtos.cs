namespace Bookstore.Domain.Addresses
{
    public class CreateAddressDto 
    { 
        public CreateAddressDto(
        string? AddressLine1,
        string? AddressLine2,
        string? City,
        string? State,
        string? Country,
        string? ZipCode,
        string? CustomerSub)
        {
            this.AddressLine1 = AddressLine1 ?? string.Empty;
            this.AddressLine2 = AddressLine2 ?? string.Empty;
            this.City = City ?? string.Empty;
            this.State = State ?? string.Empty;
            this.Country = Country ?? string.Empty;
            this.ZipCode = ZipCode ?? string.Empty;
            this.CustomerSub = CustomerSub ?? string.Empty;
        }

        public string AddressLine1 { get; }
        public string AddressLine2 { get; }
        public string City { get; }
        public string State { get; }
        public string Country { get; }
        public string ZipCode { get; }
        public string CustomerSub { get; }
    }

    public class UpdateAddressDto
    {
        public UpdateAddressDto(
        int AddressId,
        string? AddressLine1,
        string? AddressLine2,
        string? City,
        string? State,
        string? Country,
        string? ZipCode,
        string? CustomerSub)
        {
            this.AddressId = AddressId;
            this.AddressLine1 = AddressLine1 ?? string.Empty;
            this.AddressLine2 = AddressLine2 ?? string.Empty;
            this.City = City ?? string.Empty;
            this.State = State ?? string.Empty;
            this.Country = Country ?? string.Empty;
            this.ZipCode = ZipCode ?? string.Empty;
            this.CustomerSub = CustomerSub ?? string.Empty;
        }

        public int AddressId { get; }
        public string AddressLine1 { get; }
        public string AddressLine2 { get; }
        public string City { get; }
        public string State { get; }
        public string Country { get; }
        public string ZipCode { get; }
        public string CustomerSub { get; }
    }

    public class DeleteAddressDto
    {
        public DeleteAddressDto(
        int AddressId,
        string? CustomerSub)
        {
            this.AddressId = AddressId;
            this.CustomerSub = CustomerSub ?? string.Empty;
        }

        public int AddressId { get; }
        public string CustomerSub { get; }
    }
}
