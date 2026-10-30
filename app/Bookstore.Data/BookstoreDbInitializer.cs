using Bookstore.Domain.Books;
using Bookstore.Domain.ReferenceData;
using Microsoft.EntityFrameworkCore;
using System.Threading;
using System.Threading.Tasks;

namespace Bookstore.Data
{
    /// <summary>
    /// Replaces the EF6 drop-create-if-model-changes initializer that was registered through the EF6
    /// initializer API. EF Core has no initializer concept, so this class has two parts:
    ///
    ///  * <see cref="Seed"/> declares the reference data and the book catalogue as EF Core model data
    ///    (<c>HasData</c>). EF Core emits the correct <c>SET IDENTITY_INSERT</c> statements for the
    ///    explicit Id values, which a runtime <c>AddRange</c> + <c>SaveChanges</c> could not do.
    ///  * <see cref="Initialize"/> / <see cref="InitializeAsync"/> create the schema (and therefore
    ///    insert the seed data). The host must call one of them at startup, e.g.
    ///    <code>
    ///    using (var scope = app.Services.CreateScope())
    ///    {
    ///        var db = scope.ServiceProvider.GetRequiredService&lt;ApplicationDbContext&gt;();
    ///        new BookstoreDbInitializer().Initialize(db);
    ///    }
    ///    </code>
    /// </summary>
    public class BookstoreDbInitializer
    {
        /// <summary>
        /// Creates the database and schema if it does not already exist, which also inserts the data
        /// declared by <see cref="Seed"/>.
        ///
        /// Behaviour note: the EF6 initializer dropped and recreated the database whenever the model
        /// changed. EF Core has no equivalent, so this is a create-if-absent operation intended for
        /// development and for a first deployment against an empty database. It does NOT reconcile an
        /// existing database whose schema is stale — for ongoing schema changes use EF Core migrations
        /// (<c>dotnet ef migrations add</c> then <c>context.Database.Migrate()</c> in place of this call).
        /// </summary>
        /// <returns><c>true</c> if the database was created (and therefore seeded) by this call.</returns>
        public bool Initialize(ApplicationDbContext context)
        {
            return context.Database.EnsureCreated();
        }

        /// <summary>
        /// Asynchronous counterpart of <see cref="Initialize"/>.
        /// </summary>
        public async Task<bool> InitializeAsync(ApplicationDbContext context, CancellationToken cancellationToken = default)
        {
            return await context.Database.EnsureCreatedAsync(cancellationToken);
        }

        /// <summary>
        /// Declares the seed data for the model. Called from
        /// <see cref="ApplicationDbContext.OnModelCreating"/>.
        /// </summary>
        public static void Seed(ModelBuilder modelBuilder)
        {
            modelBuilder.Entity<ReferenceDataItem>().HasData(
                new ReferenceDataItem(ReferenceDataType.BookType, "Hardcover") { Id = 1 },
                new ReferenceDataItem(ReferenceDataType.BookType, "Trade Paperback") { Id = 2 },
                new ReferenceDataItem(ReferenceDataType.BookType, "Mass Market Paperback") { Id = 3 },

                new ReferenceDataItem(ReferenceDataType.Condition, "New") { Id = 4 },
                new ReferenceDataItem(ReferenceDataType.Condition, "Like New") { Id = 5 },
                new ReferenceDataItem(ReferenceDataType.Condition, "Good") { Id = 6 },
                new ReferenceDataItem(ReferenceDataType.Condition, "Acceptable") { Id = 7 },

                new ReferenceDataItem(ReferenceDataType.Genre, "Biographies") { Id = 8 },
                new ReferenceDataItem(ReferenceDataType.Genre, "Children's Books") { Id = 9 },
                new ReferenceDataItem(ReferenceDataType.Genre, "History") { Id = 10 },
                new ReferenceDataItem(ReferenceDataType.Genre, "Literature & Fiction") { Id = 11 },
                new ReferenceDataItem(ReferenceDataType.Genre, "Mystery, Thriller & Suspense") { Id = 12 },
                new ReferenceDataItem(ReferenceDataType.Genre, "Science Fiction & Fantasy") { Id = 13 },
                new ReferenceDataItem(ReferenceDataType.Genre, "Travel") { Id = 14 },

                new ReferenceDataItem(ReferenceDataType.Publisher, "Arcadia Books") { Id = 15 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Astral Publishing") { Id = 16 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Moonlight Publishing") { Id = 17 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Dreamscape Press") { Id = 18 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Enchanted Library") { Id = 19 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Fantasia House") { Id = 20 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Horizon Books") { Id = 21 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Infinity Press") { Id = 22 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Paradigm Publishing") { Id = 23 },
                new ReferenceDataItem(ReferenceDataType.Publisher, "Aurora Publishing") { Id = 24 }
            );

            modelBuilder.Entity<Book>().HasData(
                new Book("2020: The Apocalypse", "Li Juan", "6556784356", 15, 1, 13, 5, 10.95M, 25, null, null, "/Content/Images/coverimages/apocalypse.png") { Id = 1 },
                new Book("Children Of Iron", "Nikki Wolf", "7665438976", 16, 1, 11, 6, 13.95M, 3, null, null, "/Content/Images/coverimages/childrenofiron.png") { Id = 2 },
                new Book("Gold In The Dark", "Richard Roe", "5442280765", 17, 1, 13, 5, 6.50M, 10, null, null, "/Content/Images/coverimages/goldinthedark.png") { Id = 3 },
                new Book("Leagues Of Smoke", "Pat Candella", "4556789542", 18, 2, 11, 7, 3M, 1, null, null, "/Content/Images/coverimages/leaguesofsmoke.png") { Id = 4 },
                new Book("Alone With The Stars", "Carlos Salazar", "4563358087", 19, 2, 12, 5, 15.95M, 5, null, null, "/Content/Images/coverimages/alonewiththestars.png") { Id = 5 },
                new Book("The Girl In The Polaroid", "Terri Whitlock", "2354435678", 20, 1, 12, 6, 8.25M, 2, null, null, "/Content/Images/coverimages/girlinthepolaroid.png") { Id = 6 },
                new Book("1001 Jokes", "Mary Major", "6554789632", 21, 2, 11, 5, 13.95M, 7, null, null, "/Content/Images/coverimages/1001jokes.png") { Id = 7 },
                new Book("My Search For Meaning", "Mateo Jackson", "4558786554", 22, 3, 8, 7, 5M, 15, null, null, "/Content/Images/coverimages/mysearchformeaning.png") { Id = 8 }
            );
        }
    }
}
