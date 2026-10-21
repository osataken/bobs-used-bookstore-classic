using Microsoft.EntityFrameworkCore;
using System.Threading.Tasks;

namespace Bookstore.Data
{
    /// <summary>
    /// Provides database initialization helpers for EF Core.
    /// Seed data is managed via HasData() in ApplicationDbContext.OnModelCreating.
    /// </summary>
    public static class BookstoreDbInitializer
    {
        /// <summary>
        /// Applies any pending migrations and creates the database if it does not exist.
        /// Using MigrateAsync() (instead of EnsureCreatedAsync) means schema changes are
        /// tracked and applied incrementally via the migrations history table, which is
        /// required for any production schema change management.
        /// </summary>
        public static async Task InitializeAsync(ApplicationDbContext context)
        {
            await context.Database.MigrateAsync();
        }
    }
}
