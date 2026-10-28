using Microsoft.EntityFrameworkCore;
using System;

namespace Bookstore.Data
{
    /// <summary>
    /// Replacement for the EF6 <c>DropCreateDatabaseIfModelChanges</c> initializer that was
    /// registered through <c>Database.SetInitializer</c>. EF Core has no initializer concept, so
    /// the host must call <see cref="Initialize(ApplicationDbContext)"/> once at startup (see
    /// migration-summary.md). It creates the database and schema when missing; the reference data
    /// and book rows that the old Seed method inserted are applied by EF Core as part of that
    /// schema creation (see ApplicationDbContext.Seed.cs).
    /// </summary>
    public static class BookstoreDbInitializer
    {
        /// <summary>
        /// Creates the database and schema when it does not yet exist, applying the seed data
        /// configured in <c>ApplicationDbContext.Seed.cs</c>.
        /// </summary>
        /// <remarks>
        /// Behavior difference from the EF6 <c>DropCreateDatabaseIfModelChanges</c> initializer:
        /// <c>EnsureCreated</c> is a no-op when a database already exists, so it does NOT recreate
        /// the schema after a model change. When the model evolves, either drop the database so it
        /// is recreated on the next startup, or add EF Core migrations
        /// (<c>dotnet ef migrations add</c>) and call <c>context.Database.Migrate()</c> here
        /// instead — which is the recommended approach for environments with existing data.
        /// </remarks>
        public static void Initialize(ApplicationDbContext context)
        {
            if (context == null) throw new ArgumentNullException(nameof(context));

            context.Database.EnsureCreated();
        }

        public static void Initialize(string connectionString)
        {
            using (var context = new ApplicationDbContext(connectionString))
            {
                Initialize(context);
            }
        }
    }
}
