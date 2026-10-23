using Microsoft.EntityFrameworkCore;

namespace Bookstore.Data
{
    /// <summary>
    /// Creates and seeds the Bookstore database.
    ///
    /// The EF6 version of this class derived from
    /// DropCreateDatabaseIfModelChanges&lt;ApplicationDbContext&gt; and was wired up through
    /// Database.SetInitializer. EF Core has no initializer concept, so initialization is now
    /// explicit: the host calls Initialize (or InitializeAsync) once during start up, which
    /// creates the schema if it does not exist. The reference data and book seed rows are
    /// declared with HasData in ApplicationDbContext.PopulateDatabase and are written by
    /// EF Core as part of that schema creation (EF Core emits the correct
    /// SET IDENTITY_INSERT ON/OFF statements for the explicit Id values).
    ///
    /// Both methods are idempotent: when the database already exists nothing is created and
    /// no seed rows are re-inserted, so they are safe to call on every application start.
    ///
    /// Deliberate behaviour change: the EF6 initializer DROPPED AND RECREATED the database
    /// whenever the model hash changed. That is intentionally not reproduced — destroying
    /// live data on a model change is unsafe outside a throwaway development database.
    /// EnsureCreated only creates the schema when it is absent, so a later model change must
    /// be rolled out with EF Core migrations (dotnet ef migrations add / database update,
    /// which also applies the HasData seed rows) or by recreating the database by hand.
    /// </summary>
    public static class BookstoreDbInitializer
    {
        public static void Initialize(ApplicationDbContext context)
        {
            context.Database.EnsureCreated();
        }

        public static System.Threading.Tasks.Task InitializeAsync(ApplicationDbContext context)
        {
            return context.Database.EnsureCreatedAsync();
        }
    }
}
