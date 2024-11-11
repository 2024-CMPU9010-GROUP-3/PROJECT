"use client";

const TermsAndPrivacyPage = () => {
  return (
    <div className="min-h-screen bg-white text-black flex flex-col">
      <div className="bg-black text-white flex items-center justify-center h-64">
        <h1 className="text-4xl font-bold w-full text-center">Terms & Privacy</h1>
      </div>
      <div className="flex-grow container mx-auto p-8 max-w-4xl">
        <section className="mb-10">
          <h2 className="text-2xl font-semibold mb-4">Terms of Service</h2>
          <p className="mb-4">
            Welcome to our website. By accessing or using our site, you agree to comply with and be bound by the following terms and conditions. Please review them carefully.
          </p>
          <ul className="list-disc pl-5">
            <li>Use of the site is subject to all applicable laws and regulations.</li>
            <li>You agree not to use the site for any unlawful purpose.</li>
            <li>We reserve the right to modify these terms at any time.</li>
          </ul>
          <br />
          <p className="mb-4">
            We may collect the following personal data:
          </p>
          <ul className="list-disc pl-5 mb-4">
            <li>Name</li>
            <li>Email address</li>
            <li>Location data</li>
          </ul>
        </section>

        <section className="mb-10">
          <h2 className="text-2xl font-semibold mb-4">Privacy Policy</h2>
          <p className="mb-4">
            We are committed to protecting your privacy. This policy explains how we collect, use, and disclose your personal information.
          </p>
          <ul className="list-disc pl-5">
            <li>We collect information you provide directly to us.</li>
            <li>We use cookies and similar technologies to enhance your experience.</li>
            <li>We do not share your personal information with third parties without your consent, except as required by law.</li>
          </ul>
        </section>

        <section className="mb-10">
          <h2 className="text-2xl font-semibold mb-4">Your Rights</h2>
          <p className="mb-4">
            Under GDPR, you have the right to access, rectify, or erase your personal data. You also have the right to restrict or object to our processing of your data.
          </p>
          <p className="mb-4">
            If you have any questions or concerns about our terms or privacy practices, please contact us at [your contact information].
          </p>
        </section>
      </div>
    </div>
  );
};

export default TermsAndPrivacyPage;
