import { Tour } from "onborda/dist/types";
import { Eye, EyeOff } from 'lucide-react';

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
                    The current prototype works in <strong>Dublin and the immediate surroundings.</strong>
                    <br /> <br />
                    Note, the <strong>display is locked until you finish the tour!</strong>
                    <br /> <br />
                    When you&apos;re ready, hit <strong>next!</strong>
                </>,
                selector: "#onboarding-step-1",
                side: "left-top",
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
                    Note: the <strong>larger the radius- the longer it&apos;ll take to load the data!</strong>
                </>,
                selector: "#onboarding-step-2",
                side: "left-top",
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
                side: "left-top",
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
                    You can select the amenities you want to disable by clicking these:
                    <br /> <br />
                    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                        <Eye size={20} color="#3e6e96" />
                        <span style={{ marginLeft: '5rem' }}></span> {/* Adjust the margin as needed */}
                        <EyeOff size={20} color="#3e6e96" />
                    </div>

                </>,
                selector: "#onboarding-step-4",
                side: "left-top",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },


            {
                icon: <>🗺️</>,
                title: "Selection a location",
                content: <>
                    Left or right clicking anywhere on the map will display information about the selected location.
                    <br /> <br />
                    You can zoom by <strong>scrolling</strong> with a mouse or <strong>pinching</strong> on a touchpad (or touchscreen).
                </>,
                selector: "#onboarding-step-5",
                side: "right-top",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },

            {
                icon: <>🎉</>,
                title: "That's it!",
                content: <>
                    If you ever need this tour again, hit the <strong>?</strong> button in the bottom left corner.
                    <br /> <br />
                    Enjoy using <i>Magpie</i>!
                </>,
                selector: "#onboarding-step-6",
                side: "right-bottom",
                showControls: true,
                pointerPadding: 10,
                pointerRadius: 10,
            },

        ],
    }
];