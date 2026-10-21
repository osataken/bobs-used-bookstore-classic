using Microsoft.EntityFrameworkCore;
using System.Threading.Tasks;

namespace Bookstore.Data
{
    /// <summary>
    /// Provides database initialization helpers for EF Core.
    /// Seed data is now managed via HasData() in ApplicationDbContext.OnModelCreating.
    /// </summary>
    public static class BookstoreDbInitializer
    {
        /// <summary>
        /// Ensures the database is created and all pending migrations are applied.
        /// </summary>
        public static async Task InitializeAsync(ApplicationDbContext context)
        {
            await context.Database.EnsureCreatedAsync();
        }
    }
}
