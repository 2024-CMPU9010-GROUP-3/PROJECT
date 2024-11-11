import { Tour } from "onborda/dist/types";

export const steps: Tour[] = [
    {
        tour: "general-onboarding",
        steps: [
            {
                icon: <>👋</>,
                title: "Welcome to Magpie!",
                content: <>
                    Magpie is a &nbsp;
                    <a
                        href="https://en.wikipedia.org/wiki/Geographic_information_system"
                        target="_blank"
                        rel="noopener noreferrer"
                        style={{ color: 'CornflowerBlue' }}>
                        GIS application
                    </a>
                    &nbsp; that provides an &apos;at a glance&apos; view of available public services.
                    <br /> <br />
                    This tour will guide you through the application.
                    <br /> <br />
                    The current prototype works in Dublin and the immediate surroundings.
                    <br /> <br />
                    When you&apos;re ready, hit <strong>next!</strong>
                </>,
                selector: "#onboarding-step-1",
                side: "bottom",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },

            {
                icon: <>📏</>,
                title: "Search Radius",
                content: <>
                    This slider allows you to adjust the search radius. It goes all the way from 1m to 10km.
                    <br /> <br />
                    Note, the larger the radius- the longer it&apos;ll take to load the data!
                </>,
                selector: "#onboarding-step-2",
                side: "bottom",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },

            {
                icon: <>📍</>,
                title: "Marker Data",
                content: <>
                    Information about the selected location will be displayed here!
                </>,
                selector: "#onboarding-step-3",
                side: "left",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },

            {
                icon: <>⛳</>,
                title: "Selecting Amenities",
                content: <>
                    <i>Magpie</i> supports about a <strong>dozen</strong> different amenities. You may not want to see all of them at once.
                    <br /> <br />
                    You can select the amenities you want by clicking on this box, and either:
                    <br /> <br />
                    &emsp;• Typing the elements you want
                    <br />
                    &emsp;• Using the dropdown menu

                </>,
                selector: "#onboarding-step-4",
                side: "left",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },


            {
                icon: <>🗺️</>,
                title: "Selection a location",
                content: <>
                    Left or right clicking anywhere on the map will display information about the selected location.
                </>,
                selector: "#onboarding-step-5",
                side: "right",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },

            {
                icon: <>🎉</>,
                title: "That's it!",
                content: <>
                    If you ever need this tour again, hit the <strong>?</strong> button in the bottom right corner.
                    <br /> <br />
                    Enjoy using <i>Magpie</i>!
                </>,
                selector: "#onboarding-step-1",
                side: "bottom",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },

        ],
    }
];