# Minutes of the Meeting (MoM) - Week 7 - Day 3

**Project Title:** Magpie - Services at a Glance

**Group Members:** Kaustubh, Saul, Andreas, Jessica, Anais, Yuanshuo Du (Steven)

**Date:** 30th October 2024

**Time:** 12:32 PM

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

- Updates on current tasks and blockers.
- Discuss new data source integration and CSV conversion.
- Plan next steps for database and layer management.

---

### **3. Discussion and Key Points**

1. **Task Updates:**
   - **Saul:** Automated repository structure documentation and began cleaning up GitHub. Plans to rename "group three" references to "magpie" and reorganize the repository for consistency.
   - **Andreas:** Continues working on unit tests, facing delays but progressing with each component.
   - **Jessica:** Searched for new datasets on DataGov.ie, began adding intelligence to the model for larger bounding boxes, and worked on identifying empty parking spaces.
   - **Kaustubh:** Identified issues in backend data handling, which may require implementing data streaming to manage large datasets effectively. Plans to troubleshoot and update with further findings.
   - **Anais:** Focused on model tuning, especially to reduce false positives. Also worked on updating the survey results and aims to find more training images to improve accuracy.
   - **Steven:** Developed a Figma-based UI prototype and shared it with the team for input. He will continue refining the prototype and add new front-end features based on team feedback.

2. **Additional Data Sources:**
   - **Jessica and Anais** were assigned the task of converting GeoJSON datasets to CSV format for database integration. These datasets, such as parking zone maps, are intended to expand the map's functionality.
   - Saul emphasized the importance of structuring these files in the correct format, removing unnecessary columns, and adapting them to the system’s database structure.

3. **Database and Layer Management:**
   - **Kaustubh** will implement layered map views to display different data types, such as parking zones or EV charging points, allowing users to toggle layers on and off.
   - **Andreas** will create migration solutions to facilitate easier integration of new data types into the database.

4. **Data Streaming and Backend Solutions:**
   - The team discussed adding data streaming to address issues with map data overloading the front end. Andreas suggested using WebSockets to manage data flow, but the team will assess its feasibility over the coming days.

5. **Action Items:**
   - **Jessica and Anais:** Prepare additional data source CSVs, removing extraneous data columns as necessary.
   - **Andreas and Kaustubh:** Focus on layer implementation and data streaming options, planning a follow-up discussion at 7:00 PM.
   - **Steven:** Continue front-end prototyping and provide additional UI ideas for feedback.

---

**Date:** 30th October 2024
