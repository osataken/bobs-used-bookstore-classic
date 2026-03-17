using System.Collections;
using System.Collections.Generic;
using Microsoft.AspNetCore.Mvc;

namespace Bookstore.Web.ViewModel.Address
{
    public class AddressCreateUpdateViewModel
    {
        public AddressCreateUpdateViewModel() { }

        public AddressCreateUpdateViewModel(string returnUrl)
        {
            ReturnUrl = returnUrl;
        }

        public AddressCreateUpdateViewModel(Domain.Addresses.Address address, string returnUrl)
        {
            Id = address.Id;
            AddressLine1 = address.AddressLine1;
            AddressLine2 = address.AddressLine2;
            City = address.City;
            Country = address.Country;
            State = address.State;
            ZipCode = address.ZipCode;
            ReturnUrl = returnUrl;
        }

        public int Id { get; set; }

        public string AddressLine1 { get; set; }

        public string AddressLine2 { get; set; }

        public string City { get; set; }

        public string State { get; set; }

        public string Country { get; set; }

        public string ZipCode { get; set; }

        public string ReturnUrl { get; set; }

        public IEnumerable<Microsoft.AspNetCore.Mvc.Rendering.SelectListItem> States => new List<Microsoft.AspNetCore.Mvc.Rendering.SelectListItem>
        {
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "", Text = "" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "AL", Text = "Alabama" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "AK", Text = "Alaska" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "AZ", Text = "Arizona" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "AR", Text = "Arkansas" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "CA", Text = "California" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "CO", Text = "Colorado" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "CT", Text = "Connecticut" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "DE", Text = "Delaware" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "DC", Text = "District Of Columbia" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "FL", Text = "Florida" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "GA", Text = "Georgia" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "HI", Text = "Hawaii" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "ID", Text = "Idaho" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "IL", Text = "Illinois" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "IN", Text = "Indiana" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "IA", Text = "Iowa" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "KS", Text = "Kansas" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "KY", Text = "Kentucky" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "LA", Text = "Louisiana" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "ME", Text = "Maine" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "MD", Text = "Maryland" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "MA", Text = "Massachusetts" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "MI", Text = "Michigan" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "MN", Text = "Minnesota" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "MS", Text = "Mississippi" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "MO", Text = "Missouri" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "MT", Text = "Montana" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "NE", Text = "Nebraska" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "NV", Text = "Nevada" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "NH", Text = "New Hampshire" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "NJ", Text = "New Jersey" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "NM", Text = "New Mexico" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "NY", Text = "New York" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "NC", Text = "North Carolina" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "ND", Text = "North Dakota" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "OH", Text = "Ohio" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "OK", Text = "Oklahoma" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "OR", Text = "Oregon" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "PA", Text = "Pennsylvania" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "RI", Text = "Rhode Island" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "SC", Text = "South Carolina" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "SD", Text = "South Dakota" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "TN", Text = "Tennessee" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "TX", Text = "Texas" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "UT", Text = "Utah" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "VT", Text = "Vermont" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "VA", Text = "Virginia" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "WA", Text = "Washington" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "WV", Text = "West Virginia" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "WI", Text = "Wisconsin" },
            new Microsoft.AspNetCore.Mvc.Rendering.SelectListItem { Value = "WY", Text = "Wyoming" }
        };
    }

}
