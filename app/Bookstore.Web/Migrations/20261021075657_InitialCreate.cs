using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

#pragma warning disable CA1814 // Prefer jagged arrays over multidimensional

namespace Bookstore.Web.Migrations
{
    /// <inheritdoc />
    public partial class InitialCreate : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "Customer",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    Sub = table.Column<string>(type: "nvarchar(450)", maxLength: 450, nullable: false),
                    Username = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    FirstName = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    LastName = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    Email = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    DateOfBirth = table.Column<DateTime>(type: "datetime2", nullable: true),
                    Phone = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Customer", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "ReferenceData",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    DataType = table.Column<int>(type: "int", nullable: false),
                    Text = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_ReferenceData", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "ShoppingCart",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    CorrelationId = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_ShoppingCart", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "Address",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    AddressLine1 = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    AddressLine2 = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    City = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    State = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    Country = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    ZipCode = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CustomerId = table.Column<int>(type: "int", nullable: false),
                    IsActive = table.Column<bool>(type: "bit", nullable: false),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Address", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Address_Customer_CustomerId",
                        column: x => x.CustomerId,
                        principalTable: "Customer",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "Book",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    Name = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    Author = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    Year = table.Column<int>(type: "int", nullable: true),
                    ISBN = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    PublisherId = table.Column<int>(type: "int", nullable: false),
                    BookTypeId = table.Column<int>(type: "int", nullable: false),
                    GenreId = table.Column<int>(type: "int", nullable: false),
                    ConditionId = table.Column<int>(type: "int", nullable: false),
                    CoverImageUrl = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    Summary = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    Price = table.Column<decimal>(type: "decimal(18,2)", nullable: false),
                    Quantity = table.Column<int>(type: "int", nullable: false),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Book", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Book_ReferenceData_BookTypeId",
                        column: x => x.BookTypeId,
                        principalTable: "ReferenceData",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Book_ReferenceData_ConditionId",
                        column: x => x.ConditionId,
                        principalTable: "ReferenceData",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Book_ReferenceData_GenreId",
                        column: x => x.GenreId,
                        principalTable: "ReferenceData",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Book_ReferenceData_PublisherId",
                        column: x => x.PublisherId,
                        principalTable: "ReferenceData",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "Offer",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    Author = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    ISBN = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    BookName = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    FrontUrl = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    GenreId = table.Column<int>(type: "int", nullable: false),
                    ConditionId = table.Column<int>(type: "int", nullable: false),
                    PublisherId = table.Column<int>(type: "int", nullable: false),
                    BookTypeId = table.Column<int>(type: "int", nullable: false),
                    Summary = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    OfferStatus = table.Column<int>(type: "int", nullable: false),
                    Comment = table.Column<string>(type: "nvarchar(max)", nullable: true),
                    CustomerId = table.Column<int>(type: "int", nullable: false),
                    BookPrice = table.Column<decimal>(type: "decimal(18,2)", nullable: false),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Offer", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Offer_Customer_CustomerId",
                        column: x => x.CustomerId,
                        principalTable: "Customer",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_Offer_ReferenceData_BookTypeId",
                        column: x => x.BookTypeId,
                        principalTable: "ReferenceData",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Offer_ReferenceData_ConditionId",
                        column: x => x.ConditionId,
                        principalTable: "ReferenceData",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Offer_ReferenceData_GenreId",
                        column: x => x.GenreId,
                        principalTable: "ReferenceData",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Offer_ReferenceData_PublisherId",
                        column: x => x.PublisherId,
                        principalTable: "ReferenceData",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "Order",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    CustomerId = table.Column<int>(type: "int", nullable: false),
                    AddressId = table.Column<int>(type: "int", nullable: false),
                    DeliveryDate = table.Column<DateTime>(type: "datetime2", nullable: false),
                    OrderStatus = table.Column<int>(type: "int", nullable: false),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Order", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Order_Address_AddressId",
                        column: x => x.AddressId,
                        principalTable: "Address",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_Order_Customer_CustomerId",
                        column: x => x.CustomerId,
                        principalTable: "Customer",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "ShoppingCartItem",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    ShoppingCartId = table.Column<int>(type: "int", nullable: false),
                    BookId = table.Column<int>(type: "int", nullable: false),
                    Quantity = table.Column<int>(type: "int", nullable: false),
                    WantToBuy = table.Column<bool>(type: "bit", nullable: false),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_ShoppingCartItem", x => new { x.Id, x.ShoppingCartId });
                    table.ForeignKey(
                        name: "FK_ShoppingCartItem_Book_BookId",
                        column: x => x.BookId,
                        principalTable: "Book",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_ShoppingCartItem_ShoppingCart_ShoppingCartId",
                        column: x => x.ShoppingCartId,
                        principalTable: "ShoppingCart",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "OrderItem",
                columns: table => new
                {
                    Id = table.Column<int>(type: "int", nullable: false)
                        .Annotation("SqlServer:Identity", "1, 1"),
                    OrderId = table.Column<int>(type: "int", nullable: false),
                    BookId = table.Column<int>(type: "int", nullable: false),
                    Quantity = table.Column<int>(type: "int", nullable: false),
                    CreatedBy = table.Column<string>(type: "nvarchar(max)", nullable: false),
                    CreatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    UpdatedOn = table.Column<DateTime>(type: "datetime2", nullable: false),
                    RowVersion = table.Column<byte[]>(type: "rowversion", rowVersion: true, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_OrderItem", x => x.Id);
                    table.ForeignKey(
                        name: "FK_OrderItem_Book_BookId",
                        column: x => x.BookId,
                        principalTable: "Book",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_OrderItem_Order_OrderId",
                        column: x => x.OrderId,
                        principalTable: "Order",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.InsertData(
                table: "ReferenceData",
                columns: new[] { "Id", "CreatedBy", "CreatedOn", "DataType", "Text", "UpdatedOn" },
                values: new object[,]
                {
                    { 1, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(3400), 2, "Hardcover", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(3400) },
                    { 2, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4430), 2, "Trade Paperback", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4431) },
                    { 3, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4432), 2, "Mass Market Paperback", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4433) },
                    { 4, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4434), 1, "New", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4434) },
                    { 5, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4435), 1, "Like New", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4436) },
                    { 6, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4437), 1, "Good", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4437) },
                    { 7, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4438), 1, "Acceptable", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4439) },
                    { 8, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4440), 3, "Biographies", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4440) },
                    { 9, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4441), 3, "Children's Books", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4442) },
                    { 10, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4443), 3, "History", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4443) },
                    { 11, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4444), 3, "Literature & Fiction", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4445) },
                    { 12, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4446), 3, "Mystery, Thriller & Suspense", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4446) },
                    { 13, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4447), 3, "Science Fiction & Fantasy", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4448) },
                    { 14, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4449), 3, "Travel", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4449) },
                    { 15, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4450), 0, "Arcadia Books", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4450) },
                    { 16, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4451), 0, "Astral Publishing", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4452) },
                    { 17, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4453), 0, "Moonlight Publishing", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4453) },
                    { 18, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4454), 0, "Dreamscape Press", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4455) },
                    { 19, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4456), 0, "Enchanted Library", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4456) },
                    { 20, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4457), 0, "Fantasia House", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4458) },
                    { 21, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4459), 0, "Horizon Books", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4459) },
                    { 22, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4460), 0, "Infinity Press", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4460) },
                    { 23, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4461), 0, "Paradigm Publishing", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4462) },
                    { 24, "System", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4463), 0, "Aurora Publishing", new DateTime(2026, 10, 21, 7, 56, 54, 864, DateTimeKind.Utc).AddTicks(4463) }
                });

            migrationBuilder.InsertData(
                table: "Book",
                columns: new[] { "Id", "Author", "BookTypeId", "ConditionId", "CoverImageUrl", "CreatedBy", "CreatedOn", "GenreId", "ISBN", "Name", "Price", "PublisherId", "Quantity", "Summary", "UpdatedOn", "Year" },
                values: new object[,]
                {
                    { 1, "Li Juan", 1, 5, "/Content/Images/coverimages/apocalypse.png", "System", new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(1199), 13, "6556784356", "2020: The Apocalypse", 10.95m, 15, 25, null, new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(1200), null },
                    { 2, "Nikki Wolf", 1, 6, "/Content/Images/coverimages/childrenofiron.png", "System", new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(5975), 11, "7665438976", "Children Of Iron", 13.95m, 16, 3, null, new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(5976), null },
                    { 3, "Richard Roe", 1, 5, "/Content/Images/coverimages/goldinthedark.png", "System", new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(5979), 13, "5442280765", "Gold In The Dark", 6.50m, 17, 10, null, new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(5980), null },
                    { 4, "Pat Candella", 2, 7, "/Content/Images/coverimages/leaguesofsmoke.png", "System", new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6012), 11, "4556789542", "Leagues Of Smoke", 3m, 18, 1, null, new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6013), null },
                    { 5, "Carlos Salazar", 2, 5, "/Content/Images/coverimages/alonewiththestars.png", "System", new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6015), 12, "4563358087", "Alone With The Stars", 15.95m, 19, 5, null, new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6015), null },
                    { 6, "Terri Whitlock", 1, 6, "/Content/Images/coverimages/girlinthepolaroid.png", "System", new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6017), 12, "2354435678", "The Girl In The Polaroid", 8.25m, 20, 2, null, new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6018), null },
                    { 7, "Mary Major", 2, 5, "/Content/Images/coverimages/1001jokes.png", "System", new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6019), 11, "6554789632", "1001 Jokes", 13.95m, 21, 7, null, new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6020), null },
                    { 8, "Mateo Jackson", 3, 7, "/Content/Images/coverimages/mysearchformeaning.png", "System", new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6021), 8, "4558786554", "My Search For Meaning", 5m, 22, 15, null, new DateTime(2026, 10, 21, 7, 56, 54, 866, DateTimeKind.Utc).AddTicks(6022), null }
                });

            migrationBuilder.CreateIndex(
                name: "IX_Address_CustomerId",
                table: "Address",
                column: "CustomerId");

            migrationBuilder.CreateIndex(
                name: "IX_Book_BookTypeId",
                table: "Book",
                column: "BookTypeId");

            migrationBuilder.CreateIndex(
                name: "IX_Book_ConditionId",
                table: "Book",
                column: "ConditionId");

            migrationBuilder.CreateIndex(
                name: "IX_Book_GenreId",
                table: "Book",
                column: "GenreId");

            migrationBuilder.CreateIndex(
                name: "IX_Book_PublisherId",
                table: "Book",
                column: "PublisherId");

            migrationBuilder.CreateIndex(
                name: "IX_Customer_Sub",
                table: "Customer",
                column: "Sub",
                unique: true);

            migrationBuilder.CreateIndex(
                name: "IX_Offer_BookTypeId",
                table: "Offer",
                column: "BookTypeId");

            migrationBuilder.CreateIndex(
                name: "IX_Offer_ConditionId",
                table: "Offer",
                column: "ConditionId");

            migrationBuilder.CreateIndex(
                name: "IX_Offer_CustomerId",
                table: "Offer",
                column: "CustomerId");

            migrationBuilder.CreateIndex(
                name: "IX_Offer_GenreId",
                table: "Offer",
                column: "GenreId");

            migrationBuilder.CreateIndex(
                name: "IX_Offer_PublisherId",
                table: "Offer",
                column: "PublisherId");

            migrationBuilder.CreateIndex(
                name: "IX_Order_AddressId",
                table: "Order",
                column: "AddressId");

            migrationBuilder.CreateIndex(
                name: "IX_Order_CustomerId",
                table: "Order",
                column: "CustomerId");

            migrationBuilder.CreateIndex(
                name: "IX_OrderItem_BookId",
                table: "OrderItem",
                column: "BookId");

            migrationBuilder.CreateIndex(
                name: "IX_OrderItem_OrderId",
                table: "OrderItem",
                column: "OrderId");

            migrationBuilder.CreateIndex(
                name: "IX_ShoppingCartItem_BookId",
                table: "ShoppingCartItem",
                column: "BookId");

            migrationBuilder.CreateIndex(
                name: "IX_ShoppingCartItem_ShoppingCartId",
                table: "ShoppingCartItem",
                column: "ShoppingCartId");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "Offer");

            migrationBuilder.DropTable(
                name: "OrderItem");

            migrationBuilder.DropTable(
                name: "ShoppingCartItem");

            migrationBuilder.DropTable(
                name: "Order");

            migrationBuilder.DropTable(
                name: "Book");

            migrationBuilder.DropTable(
                name: "ShoppingCart");

            migrationBuilder.DropTable(
                name: "Address");

            migrationBuilder.DropTable(
                name: "ReferenceData");

            migrationBuilder.DropTable(
                name: "Customer");
        }
    }
}
