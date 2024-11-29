# Minutes of the Meeting (MoM) - Week 10 - Day 4

**Project Title:** Magpie - Services at a Glance

**Group Members:** Kaustubh, Saul, Andreas, Jessica, Anais, Yuanshuo Du (Steven)

**Date:** 21st November 2024

**Time:** 12:33 PM

**Location:** Online (Dublin, Ireland)

---

### **1. Attendance**

- **Present:**
  - Kaustubh
  - Jessica
  - Andreas
  - Anais
  - Saul
  - Steven

---

### **2. Agenda of the Meeting**

- Review UI fixes and hover feature implementation options.
- Discuss backend progress and bug fixes.
- Updates on Label Studio and parking detection evaluation.
- Summarize user testing feedback.

---

### **3. Discussion and Key Points**

1. **UI Fixes and Hover Feature Challenges:**

   - **Saul** implemented buttons for clearing screen points and toggling point visibility. However, hover functionality posed significant challenges due to limitations in React Mapbox GL.
   - Options discussed:
     - Writing a custom hover solution.
     - Refactoring the entire map component using Deck GL.
   - The team leaned towards a custom solution to avoid potential setbacks from refactoring, given the tight timeline.

2. **Backend and Onboarding Updates:**

   - **Andreas** completed a heartbeat route for backend health checks, ensuring system readiness. He will shift focus to the history feature backend, aligning it with frontend design plans.

3. **Label Studio and Parking Detection Evaluation:**

   - **Jessica** encountered discrepancies between Label Studio's exported labels and the model’s bounding box format, complicating evaluation.
   - She plans to refine the evaluation script and address format mismatches while continuing to label images for additional data.

4. **User Testing Feedback and Enhancements:**

   - **Anais** summarized feedback from user sessions, noting requests for zoom buttons on the map. Andrea also recommended making navigation more user-friendly.
   - **Steven** was tasked with implementing zoom functionality using buttons, with a branch set up for the feature.

5. **Homepage and Landing Page Progress:**

   - **Kaustubh** completed a draft for the homepage with styling updates and pushed it to the repository. He requested Saul and Jessica add machine learning and other content to the landing page.

6. **History Feature Collaboration:**

   - **Kaustubh** will coordinate with Andreas on the frontend for the history feature to align with backend development.

7. **Final Adjustments and Pull Requests:**
   - **Steven** reported that the prototype version 2 pull request is now in good shape after fixing issues. Kaustubh will review and approve it.

---

### **4. Next Steps**

- **Saul:** Work on hover functionality and assist with landing page content.
- **Jessica:** Resolve label discrepancies and refine the parking detection evaluation script.
- **Andreas:** Focus on backend integration for the history feature.
- **Anais:** Continue summarizing user testing insights and preparing for final evaluation.
- **Steven:** Implement zoom functionality and refine UI features.
- **Kaustubh:** Oversee homepage and history feature progress.

---

**Date:** 21st November 2024
