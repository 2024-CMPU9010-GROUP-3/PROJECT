# What is Magpie

Hi everyone, thank you for joining us as we demonstrate our new tool, Magpie.
Magpie is a web application that allows civil planners, and Non-Expert Users, to
gather data on public amenities in a given area. This data can be used to inform
decisions on where to build new amenities, or, to identify areas that are
*lacking* in certain amenities.

# Why is Magpie better than our competitors

Magpie is a simple, integrated solution that offers users quick access to data, at a glance. While our competitors offer similar functionality, they often require users to navigate through multiple pages, and click several buttons to access the same information. In contrast, Magpie is designed to be intuitive and easy to use.

# User personas

This is Michael O'Brian, an Urban Planning Specialist at Dublin City Council. Michael is responsible for identifying areas in Dublin that are lacking in public amenities. He uses Magpie to quickly gather data on public services, and to identify areas that are in need of additional services.

Michael has access to a wide range of data, but gaining a high-level overview of that data
can be difficult and time consuming with his current tools. Additionally, Michael needs accurate information about on-the-street parking spaces in his area, something his current setup does not provide.

# Signup page

This is the signup page, where users can create a new account. In the interest of time, we will skip
over this step and login with an existing account.

# Onboarding

On first login, Magpie provides a brief onboarding experience to help users get started. This includes
a short overview of the system, and a quick tutorial on how to use the application.

In the interest of time, we will not go through the full experience- and exit the onboarding process.

# GDPR Compliance (Accept Cookies)

Now that we've exited the onboarding process, we can see a cookie consent banner. Accepting this banner
allows Magpie to store cookies on your device. In all aspects, Magpie is GDPR compliant, and we take
the privacy of our users very seriously. We never collect user data without consent, and all location
data is processed on device, rather then on our servers.

# Tooltips

Magpie is designed to be simple and easy to use. The interface is divided into two
main sections: the map view, and the dashboard.

Hovering over buttons will display tooltips, providing users with additional information on what each button does.

# Discuss UI elements

The map view displays information about the current area, with buttons in the top left and bottom left to indicate special functions. In the top, we have buttons to: logout, view saved locations, and restart the onboarding process. In the bottom, we have buttons to zoom if the user is using a device without a scroll wheel.

# Search (grangorman)

To the right, we have the dashboard. The dashboard begins with search options, this card allows users to select their search radius, either by using a slider, or by selecting from a set of presets. Below this, we have a search box, typing into this box and selecting the search button, will perform a location search.

# Zoom in using scroll wheel

Now that we have a marker, we can zoom in using the scroll wheel on our mouse, to gain a better perspective of the retrieved amenities. Note that, amenities are clustered together, and zooming in will reveal more detailed information.

# Set radius to 500m

Our search radius is currently set to it's default of 4 kilometres, we can change this by selecting one of the search presets. For this demonstration, we will set the radius to 500 meters.

# Amenity Counts

Looking at the map, we can see a number of amenities. By looking to the right, at the dashboard, we can see a breakdown of their numbers, as well as their types.

# Amenity Filters

By clicking on the blue 'eye' buttons, we can filter out amenities.

## Disable all

The top eye icon, will disable all amenities on the map, giving us a better view of the geography.

## Enable all

We can press the button again, to re-enable all of the amenities.

## Disable parking

As you can see, there are many parking spaces in the area. We can disable just the parking amenity by selecting the 'eye' next to parking.

## Reenable parking

We can re-enable parking by selecting the 'eye' again.

# Hover over (Bike Stand)

Hovering over an amenity will display a tooltip, providing additional information about that amenity.

# Hover over Parking, discuss ML

Parking spaces are detected by our machine learning model, which utilises satellite images to identify public and private parking spaces. In the tooltip, we can see that the parking space is public, and is in the red zone- indicating that this parking space is in high demand, as defined by Dublin City Council.

# Zoom out using Buttons

Now that we have explored this area, we can zoom out using the buttons in the bottom left of the map.

# Saving a search

If we find an area that we would like to return to, we can save the search by clicking the 'save' button in the dashboard. This will save the search, and allow us to return to it at a later date.

# Open up saved locations

To view saved locations, we can click the 'saved locations' button in the top left of the map.

# Discuss saved locations page

The saved locations page displays a list of all saved locations. This view, also displays the radius of the search, the amenities found in the search, when the search was performed and a button to retrieve the search.

As you can see there multiple saved locations, a civil planner could use this view to compare different areas, and make informed decisions on where to build new amenities.

If we decide a particular location is no longer relevant, we can select it, and click the 'delete' button to remove it from the list.

# Load the saved locations

By clicking the 'show on map' button, we can load the saved location onto the map. This fully retrieves the search, correctly setting the radius and amenities filter.

# Log out

Once we are finished looking at amenities, we can click the logout button to return us to the login page. This completes our demonstration of Magpie.

# Thank yous

Thank you for joining us today, we would like to thank our lecturers, our supervisor, our test users and our survey respondents for their help in developing Magpie. We hope you enjoyed our demonstration.