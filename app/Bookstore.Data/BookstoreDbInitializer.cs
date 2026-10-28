#if NET10_0_OR_GREATER
using Bookstore.Domain.Books;
using Bookstore.Domain.ReferenceData;
using Microsoft.EntityFrameworkCore;
using System.Collections.Generic;
using System.Linq;

namespace Bookstore.Data
{
    /// <summary>
    /// EF Core replacement for the EF6 <c>DropCreateDatabaseIfModelChanges</c> initializer that was
    /// previously registered through <c>Database.SetInitializer</c>. EF Core has no initializer
    /// concept, so the host application must invoke <see cref="Initialize"/> at startup (where
    /// <c>SetInitializer</c> used to run) to create the schema and seed the reference data.
    /// <para>
    /// Note: EF Core has no equivalent of "drop and recreate when the model changes", so an
    /// existing database whose schema has drifted is left in place rather than dropped. A host
    /// that wants the old development behaviour can call
    /// <c>context.Database.EnsureDeleted()</c> before <see cref="Initialize"/>; production
    /// deployments should use EF Core migrations (<c>Database.Migrate()</c>) instead.
    /// </para>
    /// </summary>
    public class BookstoreDbInitializer
    {
        /// <summary>
        /// Creates the database and schema if they do not already exist, then seeds the
        /// reference data and the initial book catalogue.
        /// </summary>
        public static void Initialize(ApplicationDbContext context)
        {
            context.Database.EnsureCreated();

            Seed(context);
        }

        public static void Seed(ApplicationDbContext context)
        {
            // Guard against re-seeding on application restart. The EF6 initializer only ran
            // Seed() when the database was (re)created, so an already-populated database must
            // be left untouched.
            if (context.ReferenceData.Any()) return;

            var referenceDataItems = new List<ReferenceDataItem> {
                new ReferenceDataItem(ReferenceDataType.BookType, "Hardcover") { Id = 1 },
                new ReferenceDataItem(ReferenceDataType.BookType, "Trade Paperback") { Id = 2 },
                new ReferenceDataItem(ReferenceDataType.BookType, "Mass Market Paperback") { Id = 3 },

                new ReferenceDataItem(ReferenceDataType.Condition, "New") { Id = 4 },
                new ReferenceDataItem(ReferenceDataType.Condition, "Like New") { Id = 5 },
                new ReferenceDataItem(ReferenceDataType.Condition, "Good") { Id = 6 },
                new ReferenceDataItem(ReferenceDataType.Condition, "Acceptable") { Id = 7 },

                new ReferenceDataItem(ReferenceDataType.Genre, "Biographies") { Id = 8 },
                new ReferenceDataItem(ReferenceDataType.Genre, "Children's Books") { Id = 9 },
                new ReferenceDataItem( ReferenceDataType.Genre, "History") { Id = 10 },
                new ReferenceDataItem( ReferenceDataType.Genre, "Literature & Fiction") { Id = 11 },
                new ReferenceDataItem( ReferenceDataType.Genre, "Mystery, Thriller & Suspense") { Id = 12 },
                new ReferenceDataItem( ReferenceDataType.Genre, "Science Fiction & Fantasy") { Id = 13 },
                new ReferenceDataItem( ReferenceDataType.Genre, "Travel") { Id = 14 },

                new ReferenceDataItem( ReferenceDataType.Publisher, "Arcadia Books") { Id = 15 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Astral Publishing") { Id = 16 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Moonlight Publishing") { Id = 17 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Dreamscape Press") { Id = 18 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Enchanted Library") { Id = 19 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Fantasia House") { Id = 20 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Horizon Books") { Id = 21 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Infinity Press") { Id = 22 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Paradigm Publishing") { Id = 23 },
                new ReferenceDataItem( ReferenceDataType.Publisher, "Aurora Publishing") { Id = 24 }
           };

            var books = new List<Book> {
                new Book("2020: The Apocalypse", "Li Juan", "6556784356", 15, 1, 13, 5, 10.95M, 25, null, null, "/Content/Images/coverimages/apocalypse.png") { Id = 1 },
                new Book("Children Of Iron", "Nikki Wolf", "7665438976", 16, 1, 11, 6, 13.95M, 3, null, null, "/Content/Images/coverimages/childrenofiron.png") { Id = 2 },
                new Book("Gold In The Dark", "Richard Roe", "5442280765", 17, 1, 13, 5, 6.50M, 10, null, null, "/Content/Images/coverimages/goldinthedark.png") { Id = 3 },
                new Book("Leagues Of Smoke", "Pat Candella", "4556789542", 18, 2, 11, 7, 3M, 1, null, null, "/Content/Images/coverimages/leaguesofsmoke.png") { Id = 4 },
                new Book("Alone With The Stars", "Carlos Salazar", "4563358087", 19, 2, 12, 5, 15.95M, 5, null, null, "/Content/Images/coverimages/alonewiththestars.png") { Id = 5 },
                new Book("The Girl In The Polaroid", "Terri Whitlock", "2354435678", 20, 1, 12, 6, 8.25M, 2, null, null, "/Content/Images/coverimages/girlinthepolaroid.png") { Id = 6 },
                new Book("1001 Jokes", "Mary Major", "6554789632", 21, 2, 11, 5, 13.95M, 7, null, null, "/Content/Images/coverimages/1001jokes.png") { Id = 7 },
                new Book("My Search For Meaning", "Mateo Jackson", "4558786554", 22, 3, 8, 7, 5M, 15, null, null, "/Content/Images/coverimages/mysearchformeaning.png") { Id = 8 }
            };

            // The seed data supplies explicit primary keys. EF Core would otherwise let the
            // store generate them, so IDENTITY_INSERT must be enabled for the duration of the
            // insert to preserve the exact Ids the reference data FKs in the book rows rely on.
            InsertWithExplicitIds(context, "ReferenceData", () => context.ReferenceData.AddRange(referenceDataItems));
            InsertWithExplicitIds(context, "Book", () => context.Book.AddRange(books));
        }

        private static void InsertWithExplicitIds(ApplicationDbContext context, string tableName, System.Action addEntities)
        {
            addEntities();

            if (context.Database.IsSqlServer())
            {
                using (var transaction = context.Database.BeginTransaction())
                {
                    // tableName is always a known constant ("ReferenceData" or "Book"), never user input.
#pragma warning disable EF1002
                    context.Database.ExecuteSqlRaw($"SET IDENTITY_INSERT [{tableName}] ON");

                    try
                    {
                        context.SaveChanges();
                    }
                    finally
                    {
                        context.Database.ExecuteSqlRaw($"SET IDENTITY_INSERT [{tableName}] OFF");
                    }
#pragma warning restore EF1002

                    transaction.Commit();
                }
            }
            else
            {
                context.SaveChanges();
            }
        }
    }
}
#endif
