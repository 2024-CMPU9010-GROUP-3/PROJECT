# Minutes of the Meeting (MoM) - Week 11 - Day 3

**Project Title:** Magpie - Services at a Glance

**Group Members:** Kaustubh, Saul, Andreas, Jessica, Anais, Yuanshuo Du (Steven)

**Date:** 27th November 2024

**Time:** 12:34 PM

**Location:** Online (Dublin, Ireland)

---

### **1. Attendance**

- **Present:**
  - Kaustubh
  - Saul
  - Andreas
  - Jessica
  - Anais
  - Steven

---

### **2. Agenda of the Meeting**

- Updates on tasks completed yesterday and plans for today.
- Discuss parking classification refinements and addressing road misclassification issues.
- Review Docker and database connection challenges.

---

### **3. Discussion and Key Points**

1. **Individual Updates:**

   - **Anais:**
     - Worked on her personal diary, sent out more emails to schedule user testing sessions, and completed retrospective documentation. Plans to continue cleaning files and sending emails today.
   - **Saul:**
     - Completed the card hover feature. Plans to review issues assigned to him and work on tasks related to suggestions received from Anais.
   - **Andreas:**
     - Resolved a long-standing backend issue by implementing graceful shutdown functionality. Shared the final report structure in the GitHub repository under the `final-report` branch. Plans to address open issues today.
   - **Jessica:**
     - Added a confidence threshold to the parking detection model and labeled parking spots as public or private. Raised questions about focusing on classification improvements, adding new categories like parking lots, and addressing highway misclassification issues.
   - **Steven:**
     - Attempted to resolve issues with unnecessary icons and menu items on the UI. Plans to discuss placement of the history page and download button with Kaustubh. Mentioned ongoing Docker-related database connection problems.

2. **Parking Classification and Model Refinements:**

   - The team discussed Jessica’s proposal to add a parking lot category and refine classifications for on-street and private parking. While it was noted that parking lots could be public or private, adding a general parking lot label was agreed upon to improve data representation.

3. **Docker and Database Connection Issues:**

   - **Steven** reported that Docker fails to connect to the database on his machine, leading to a 500 error on the front end.
   - **Andreas** suggested checking environment variables in the `.env` file and ensuring the backend is running properly within Docker.

4. **Feedback on Expanding Website Testing:**

   - The team discussed whether to expand user testing to a broader audience, such as casual users. It was decided to focus on receiving feedback from professional users, such as architects and urban planners, before considering a public release.

5. **Reminders:**
   - **Anais** reminded team members to submit retrospectives by 4 PM.
   - **Kaustubh** will finalize the homepage updates before his next availability.

---

### **4. Next Steps**

- **Anais:** Continue scheduling user testing sessions and cleaning up files.
- **Saul:** Review assigned issues and suggestions, begin implementing related fixes.
- **Andreas:** Address open backend issues and provide support for Docker troubleshooting.
- **Jessica:** Refine parking classification model and evaluate additional categories.
- **Steven:** Resolve database connection issues and finalize UI adjustments.

---

**Date:** 27th November 2024
