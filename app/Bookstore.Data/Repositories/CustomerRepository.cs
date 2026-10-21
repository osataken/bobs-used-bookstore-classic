using Bookstore.Domain.Customers;
using Microsoft.EntityFrameworkCore;
using System;
using System.Threading.Tasks;

namespace Bookstore.Data.Repositories
{
    public class CustomerRepository : ICustomerRepository
    {
        private readonly ApplicationDbContext dbContext;

        public CustomerRepository(ApplicationDbContext dbContext)
        {
            this.dbContext = dbContext;
        }

        async Task ICustomerRepository.AddAsync(Customer customer)
        {
            await dbContext.Customer.AddAsync(customer);
        }

        async Task<Customer> ICustomerRepository.GetAsync(int id)
        {
            return await dbContext.Customer.FindAsync(id) ?? throw new InvalidOperationException($"Customer with id {id} not found.");
        }

        async Task<Customer> ICustomerRepository.GetAsync(string sub)
        {
            return await dbContext.Customer.SingleOrDefaultAsync(x => x.Sub == sub) ?? throw new InvalidOperationException($"Customer with sub '{sub}' not found.");
        }

        async Task ICustomerRepository.SaveChangesAsync()
        {
            await dbContext.SaveChangesAsync();
        }
    }
}
