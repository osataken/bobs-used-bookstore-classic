namespace Bookstore.Domain.AdminUser
{
    public class DbSecrets
    {
        public string Host { get; set; } = string.Empty;

        public int Port { get; set; }

        public string Username { get; set; } = string.Empty;

        public string Password { get; set; } = string.Empty;

        public string DbInstanceIdentifier { get; set; } = string.Empty;

        public string Engine { get; set; } = string.Empty;
    }
}
