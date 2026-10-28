FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

COPY . ./

RUN dotnet restore "app/Bookstore.Web/Bookstore.Web.csproj"

RUN dotnet publish "app/Bookstore.Web/Bookstore.Web.csproj" -c Release -o /app/publish --no-restore

FROM mcr.microsoft.com/dotnet/aspnet:10.0 AS runtime

WORKDIR /app

COPY --from=build /app/publish/ .

EXPOSE 8080

ENTRYPOINT ["dotnet", "Bookstore.Web.dll"]
