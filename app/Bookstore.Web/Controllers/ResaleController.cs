using System.Threading.Tasks;
using Bookstore.Web.ViewModel.Resale;
using Bookstore.Web.Helpers;
using Bookstore.Domain.Offers;
using Bookstore.Domain.ReferenceData;
using Microsoft.AspNetCore.Mvc;

namespace Bookstore.Web.Controllers
{
    public class ResaleController : Controller
    {
        private readonly IReferenceDataService referenceDataService;
        private readonly IOfferService offerService;

        public ResaleController(IReferenceDataService referenceDataService, IOfferService offerService)
        {
            this.referenceDataService = referenceDataService;
            this.offerService = offerService;
        }

        public async Task<ActionResult> Index()
        {
            var offers = await offerService.GetOffersAsync(User.GetSub());

            return View(new ResaleIndexViewModel(offers));
        }

        public async Task<ActionResult> Create()
        {
            var referenceDataDtos = await referenceDataService.GetAllReferenceDataAsync();

            return View(new ResaleCreateViewModel(referenceDataDtos));
        }

        [HttpPost]
        public async Task<ActionResult> Create(ResaleCreateViewModel resaleViewModel)
        {
            if (!ModelState.IsValid)
            {
                // The select lists are not round-tripped by model binding, so they have to be
                // rehydrated before the Create view is re-displayed with the validation errors.
                var referenceData = await referenceDataService.GetAllReferenceDataAsync();

                return View(new ResaleCreateViewModel(referenceData)
                {
                    BookName = resaleViewModel.BookName,
                    Author = resaleViewModel.Author,
                    ISBN = resaleViewModel.ISBN,
                    BookPrice = resaleViewModel.BookPrice,
                    SelectedBookTypeId = resaleViewModel.SelectedBookTypeId,
                    SelectedConditionId = resaleViewModel.SelectedConditionId,
                    SelectedGenreId = resaleViewModel.SelectedGenreId,
                    SelectedPublisherId = resaleViewModel.SelectedPublisherId
                });
            }

            var dto = new CreateOfferDto(
                User.GetSub(), 
                resaleViewModel.BookName, 
                resaleViewModel.Author, 
                resaleViewModel.ISBN, 
                resaleViewModel.SelectedBookTypeId, 
                resaleViewModel.SelectedConditionId, 
                resaleViewModel.SelectedGenreId, 
                resaleViewModel.SelectedPublisherId, 
                resaleViewModel.BookPrice);

            await offerService.CreateOfferAsync(dto);

            return RedirectToAction(nameof(Index));
        }
    }
}