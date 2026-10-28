# Build stage: full SDK, Linux/amd64.
FROM --platform=linux/amd64 mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

# Project files first, so restore is cached independently of source changes.
# Bookstore.Web references Bookstore.Common, Bookstore.Data and Bookstore.Domain.
COPY ["app/Bookstore.Web/Bookstore.Web.csproj", "app/Bookstore.Web/"]
COPY ["app/Bookstore.Common/Bookstore.Common.csproj", "app/Bookstore.Common/"]
COPY ["app/Bookstore.Data/Bookstore.Data.csproj", "app/Bookstore.Data/"]
COPY ["app/Bookstore.Domain/Bookstore.Domain.csproj", "app/Bookstore.Domain/"]
RUN dotnet restore "app/Bookstore.Web/Bookstore.Web.csproj"

# Now the sources of the web project and its three referenced projects.
COPY ["app/Bookstore.Web/", "app/Bookstore.Web/"]
COPY ["app/Bookstore.Common/", "app/Bookstore.Common/"]
COPY ["app/Bookstore.Data/", "app/Bookstore.Data/"]
COPY ["app/Bookstore.Domain/", "app/Bookstore.Domain/"]

RUN dotnet publish "app/Bookstore.Web/Bookstore.Web.csproj" \
    -c Release \
    --no-restore \
    -o /app/publish

# Runtime stage: ASP.NET Core runtime only.
FROM --platform=linux/amd64 mcr.microsoft.com/dotnet/aspnet:10.0 AS final
WORKDIR /app

ENV ASPNETCORE_URLS=http://+:8080 \
    DOTNET_RUNNING_IN_CONTAINER=true

EXPOSE 8080

COPY --from=build /app/publish .

# The aspnet image already defines the non-root "app" user.
USER app

ENTRYPOINT ["dotnet", "Bookstore.Web.dll"]
