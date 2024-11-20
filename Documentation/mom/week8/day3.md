# Minutes of the Meeting (MoM) - Week 8 - Day 3

**Project Title:** Magpie - Services at a Glance

**Group Members:** Kaustubh, Saul, Andreas, Jessica, Anais, Yuanshuo Du (Steven)

**Date:** 6th November 2024

**Time:** 12:33 PM

**Location:** Online (Dublin, Ireland)

---

### **1. Attendance**

- **Present:**
  - Kaustubh
  - Jessica
  - Andreas
  - Steven
  - Anais
  - Saul

---

### **2. Agenda of the Meeting**

- Updates on survey analysis, usability testing, and expert review.
- Review UI/UX updates and data handling improvements.
- Discuss technical issues related to session management and environment configurations.

---

### **3. Discussion and Key Points**

1. **Survey Analysis and Usability Testing Preparations:**
   - **Anais** shared progress on survey analysis and prepared questions for usability testing, aiming to send them out by tomorrow. She clarified the testing approach, which will involve controlled and uncontrolled user tasks followed by a survey to gather feedback.
   - An expert review request will be sent to Andrea after verifying the filter functionality. Anais suggested scheduling this review for early next week to accommodate other ongoing tasks.

2. **Data Security and Environment Configuration:**
   - **Saul** is implementing SOPs (Secrets Operations Procedures) for environment variables to improve data security. This will allow encrypted files to be shared securely within the team, resolving recent issues with missing environment variables.
   - Saul explained that SOPs would streamline repository management and make it safer to upload sensitive files to GitHub.

3. **Technical Challenges with Parking Detection:**
   - **Jessica** highlighted issues with YOLO's parking detection, particularly in residential areas with varying parking orientations. Misclassifications are occurring due to overlapping bounding boxes, affecting accuracy.
   - **Jessica** plans to refine bounding box definitions and will send updated model weights to test if the changes improve detection accuracy.

4. **Session and Cookie Management Discussion:**
   - **Steven** reported on updates to session and cookie management, explaining challenges with the cookie consent banner. The team debated the best approach to align with GDPR while maintaining session functionality. Kaustubh suggested using sessions over cookies to control login states, ensuring that the system only saves sessions if consent is given.
   - The discussion covered the distinction between cookies and sessions, with the team agreeing to refactor login handling to better comply with privacy standards.

5. **UI/UX and Dashboard Updates:**
   - **Kaustubh** completed a major UI update, making the dashboard fully responsive for mobile and tablet use. This includes updated layouts and dashboard components.
   - **Steven** confirmed that privacy and terms pages are ready for final review, along with improvements to session management for a smoother user experience.

6. **Next Steps and Prototype Documentation:**
   - The team discussed documenting prototype versions in the user manual to track changes informed by user feedback. Kaustubh agreed that this approach would help tell the development story effectively.
   - **Andreas** is finalizing unit tests and preparing to work on filtering functionality. He will address any issues from testing and focus on making the UI more user-friendly based on the upcoming usability results.

---

**Date:** 6th November 2024
